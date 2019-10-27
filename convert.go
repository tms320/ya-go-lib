package yagolib

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// TryToConvert tries to convert 'src' of arbitrary type to target variable
// of arbitrary type pointed to by 'dstPtr':
// var i int
// yagolib.TryToConvert("2019", &i, nil)
// For some conversions you should pass additional parameter via 'param'.
// For example, if you convert 'string' type to 'time.Time' then you should
// set 'param' to time layout:
// var t time.Time
// yagolib.TryToConvert("2019-10-27T18:42:09+03:00", &t, time.RFC3339)
func TryToConvert(src, dstPtr, param interface{}) error {
	if reflect.TypeOf(dstPtr).Kind() == reflect.Ptr {
		dstVal := reflect.ValueOf(dstPtr).Elem()
		srcStrOrig := fmt.Sprint(src)
		srcStr := strings.Trim(srcStrOrig, ` "'`)
		var err error
		switch dstVal.Kind() {
		case reflect.Bool:
			falsesTrues := [...]string{
				"false", "off", "no", "0", "-",
				"true", "on", "yes", "1", "+"}
			srcStr = strings.ToLower(srcStr)
			for i, s := range falsesTrues {
				if srcStr == s {
					dstVal.SetBool(i >= (len(falsesTrues) >> 1))
					return nil
				}
			}
			err = fmt.Errorf(`parsing "%s": invalid syntax`, srcStrOrig)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			base := GetBaseOfIntString(srcStr)
			switch base {
			case 16:
				srcStr = strings.TrimPrefix(srcStr, "0x")
				srcStr = strings.TrimSuffix(srcStr, "h")
			case 2:
				srcStr = strings.TrimSuffix(srcStr, "b")
			}
			if v, e := strconv.ParseInt(srcStr, base, int(dstVal.Type().Size())*8); e == nil {
				dstVal.SetInt(v)
			} else {
				err = e
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			base := GetBaseOfIntString(srcStr)
			switch base {
			case 16:
				srcStr = strings.TrimPrefix(srcStr, "0x")
				srcStr = strings.TrimSuffix(srcStr, "h")
			case 2:
				srcStr = strings.TrimSuffix(srcStr, "b")
			}
			if v, e := strconv.ParseUint(srcStr, base, int(dstVal.Type().Size())*8); e == nil {
				dstVal.SetUint(v)
			} else {
				err = e
			}
		case reflect.Float32:
			if v, e := strconv.ParseFloat(srcStr, 32); e == nil {
				dstVal.SetFloat(v)
			} else {
				err = e
			}
		case reflect.Float64:
			if v, e := strconv.ParseFloat(srcStr, 64); e == nil {
				dstVal.SetFloat(v)
			} else {
				err = e
			}
		case reflect.String:
			dstVal.SetString(srcStrOrig)
		default:
			if dstVal.Type().String() == "time.Time" {
				var t time.Time
				var e error
				ok := false
				timeLayouts := []string{
					time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z,
					time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.Kitchen,
					"2006-01-02 15:04:05", "02.01.2006 15:04:05", "15:04:05 02.01.2006",
					"2006-01-02 15:04:05 MST", "2006-01-02 15:04:05 -0700",
					"2006-01-02 15:04:05 -0700 MST"}
				if param != nil {
					timeLayouts = append(timeLayouts, fmt.Sprint(param))
				}
				for _, layout := range timeLayouts {
					if t, e = time.Parse(layout, srcStrOrig); e == nil {
						ok = true
						break
					}
				}
				if !ok {
					var unixTime int64
					var unixTimeF float64
					if unixTime, e = strconv.ParseInt(srcStr, 10, 64); e == nil { // Unix time?
						t = time.Unix(unixTime, 0)
						ok = true
					} else if unixTimeF, e = strconv.ParseFloat(srcStr, 64); e == nil {
						t = time.Unix(int64(unixTimeF), 0)
						ok = true
					}
				}
				if ok {
					dstVal.Set(reflect.ValueOf(t))
					return nil
				}
				err = fmt.Errorf(`parsing "%s": unknown format`, srcStrOrig)
			} else {
				return fmt.Errorf("Target type '%s' is not supported", dstVal.Type())
			}
		}
		if err != nil {
			return fmt.Errorf("Can't convert type '%v' to '%v': %s", reflect.TypeOf(src), dstVal.Type(), err.Error())
		}
		return nil
	}
	return errors.New("'dstPtr' is not pointer")
}

// ParseMapToStruct maps 'srcMap' to target structure pointed to by 'dstPtr'.
// The value of key stored to the structure field if their names are similar.
// (the case of symbols and '-'/'_' chars are ignored)
// For example, the following pairs of names will be considered identical:
// Var_name/varName, VarName/var-name, VarName/varname.
// In addition, the structure field may have a tag to define alternative name:
// type testStruct struct {
//	  Value1   int `intVal`	// tag `intVal` is alternative name
//	  FloatVal float64
// }
// var ts testStruct
// m := map[string]interface{}{"int_val": 2019, "float_val": 20.19}
// yagolib.ParseMapToStruct(m, &ts)
// The 'ts' struct now has values: {2019, 20.19}.
// Attention! The structure fields must be exported (the first char of name must be capitalized).
// The function returns the number of successfully mapped keys and error.
func ParseMapToStruct(srcMap map[string]interface{}, dstPtr interface{}) (int, error) {
	var errMsg string
	fieldsCnt := 0
	if reflect.TypeOf(dstPtr).Kind() == reflect.Ptr {
		if reflect.TypeOf(dstPtr).Elem().Kind() == reflect.Struct {
			structValue := reflect.ValueOf(dstPtr).Elem()
			structType := structValue.Type()
			for i := 0; i < structType.NumField(); i++ { // iterate through the structure fields
				field := structType.Field(i)
				fieldValue := structValue.Field(i)
				fieldNormName := RemoveCharacters(field.Name, "-_")
				fieldNormTag := RemoveCharacters(string(field.Tag), "-_ ")
				for srcKey, srcValue := range srcMap { // search the key of the map that matches structure field
					itemNormName := RemoveCharacters(srcKey, "-_ ")
					if strings.EqualFold(itemNormName, fieldNormName) || strings.EqualFold(itemNormName, fieldNormTag) {
						if fieldValue.IsValid() {
							if fieldValue.CanSet() {
								if err := TryToConvert(srcValue, fieldValue.Addr().Interface(), nil); err == nil {
									fieldsCnt++
								} else {
									errMsg += fmt.Sprintf("Can't set field '%v': %v\n", field.Name, err.Error())
								}
							} else {
								if IsFirstRuneUpper(field.Name) {
									errMsg += fmt.Sprintf("Can't set field '%v'\n", field.Name)
								} else {
									errMsg += fmt.Sprintf("Can't set field '%v': the field is unexported\n", field.Name)
								}
							}
						} else {
							errMsg += fmt.Sprintf("Field '%v' is not valid\n", field.Name)
						}
					}
				}
			}
		} else {
			errMsg = "'dstPtr' must be pointer to structure"
		}
	} else {
		errMsg = "'dstPtr' must be pointer"
	}
	if errMsg == "" {
		return fieldsCnt, nil
	}
	return fieldsCnt, errors.New(strings.TrimRight(errMsg, "\r\n "))
}
