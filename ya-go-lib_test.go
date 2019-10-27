package yagolib

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestIsFirstRuneUpper(t *testing.T) {
	type test struct {
		in  string
		out bool
	}
	tests := [...]test{
		{"Upper", true}, {"lower", false}, {"Верхний", true}, {"нижний", false}, {"", false}}
	for _, tt := range tests {
		result := IsFirstRuneUpper(tt.in)
		if result != tt.out {
			t.Errorf(`IsFirstRuneUpper("%v") returned %v; expected: %v`, tt.in, result, tt.out)
		}
	}

}

func TestRemoveCharacters(t *testing.T) {
	source := "Test String"
	charsToRemove := "st "
	expected := "TeSring"
	result := RemoveCharacters(source, charsToRemove)
	if result != expected {
		t.Errorf("RemoveCharacters(%v, %v) returned '%v'; expected: '%v'", source, charsToRemove, result, expected)
	}
}

func TestGetBaseOfIntString(t *testing.T) {
	type test struct {
		in  string
		out int
	}
	tests := [...]test{
		{"1976", 10},
		{"0x1976", 16}, {"0X1976", 16}, {"1976h", 16}, {"1976H", 16},
		{"0b10101010", 2}, {"0B10101010", 2}, {"10101010b", 2}, {"10101010B", 2}}
	for _, tt := range tests {
		result := GetBaseOfIntString(tt.in)
		if result != tt.out {
			t.Errorf("GetBaseOfIntString(%v) returned %v; expected: %v", tt.in, result, tt.out)
		}
	}
}

func TestTryToConvert(t *testing.T) {
	type test struct {
		in1, in2, in3, out interface{}
	}
	var s string
	var i int
	var ui uint
	var b bool
	var f32 float32
	var f64 float64
	var tm time.Time
	var dstNotSupported interface{}
	tests := [...]test{
		{"'dst' not pointer", 34, nil, nil}, // out=nil - so the function being tested must return error
		{true, &b, nil, true},
		{false, &b, nil, false},
		{"true", &b, nil, true},
		{"false", &b, nil, false},
		{"on", &b, nil, true},
		{"off", &b, nil, false},
		{"yes", &b, nil, true},
		{"no", &b, nil, false},
		{"1", &b, nil, true},
		{"0", &b, nil, false},
		{"+", &b, nil, true},
		{"-", &b, nil, false},
		{"error", &b, nil, nil},
		{1976, &i, nil, 1976},
		{"1976", &i, nil, 1976},
		{-1976, &i, nil, -1976},
		{"-1976", &i, nil, -1976},
		{"0x1976", &i, nil, 0x1976},
		{"1976h", &i, nil, 0x1976},
		{"10101010b", &i, nil, 170},
		{"error", &i, nil, nil},
		{19.76, &i, nil, nil},
		{1976, &ui, nil, 1976},
		{"1976", &ui, nil, 1976},
		{"0x1976", &ui, nil, 0x1976},
		{"1976h", &ui, nil, 0x1976},
		{"10101010b", &ui, nil, 170},
		{"error", &ui, nil, nil},
		{19.76, &f32, nil, 19.76},
		{1976, &f32, nil, 1976},
		{"19.76", &f32, nil, 19.76},
		{"error", &f32, nil, nil},
		{19.76, &f64, nil, 19.76},
		{1976, &f64, nil, 1976},
		{"19.76", &f64, nil, 19.76},
		{"error", &f64, nil, nil},
		{"1976-01-03 13:32:54", &tm, nil, "1976-01-03 13:32:54 +0000 UTC"},
		{"5 Dec 1954 year 13h 32min 54sec", &tm, "2 Jan 2006 year 15h 04min 05sec", "1954-12-05 13:32:54 +0000 UTC"},
		{"189539641", &tm, nil, time.Unix(189539641, 0).String()},            // out = "1976-01-03 17:54:01 +0000 UTC"
		{"189539641.0", &tm, nil, time.Unix(int64(189539641.0), 0).String()}, // out = "1976-01-03 17:54:01 +0000 UTC"
		{"error", &tm, nil, nil},
		{"1976", &s, nil, "1976"},
		{1976, &s, nil, "1976"},
		{19.76, &s, nil, "19.76"},
		{"target type not supported", &dstNotSupported, nil, nil}}
	for _, tt := range tests {
		in1 := fmt.Sprint(tt.in1)
		if reflect.TypeOf(tt.in1).Kind() == reflect.String {
			in1 = strconv.Quote(in1)
		}
		if err := TryToConvert(tt.in1, tt.in2, tt.in3); err == nil {
			result := fmt.Sprint(reflect.ValueOf(tt.in2).Elem())
			if result != fmt.Sprint(tt.out) {
				t.Errorf("TryToConvert(%v, &dst, %v) returned dst = %v; expected: %v",
					in1, tt.in3, result, tt.out)
			}
		} else {
			if tt.out != nil {
				t.Errorf("TryToConvert(%v, &dst, %v) returned error: '%v'; expected: %v",
					in1, tt.in3, err.Error(), tt.out)
			}
		}
	}
}

func TestParseMapToStruct(t *testing.T) {
	type test struct {
		// Input data:
		srcMap map[string]interface{}
		dstPtr interface{}
		// Results expected:
		exStruct interface{}
		exNum    int
		exErr    bool // flag that the error must be returned (not 'nil')
	}
	tests := [...]test{
		// Test 0
		{
			map[string]interface{}{},
			struct{}{},
			struct{}{},
			0, true,
		},
		// Test 1
		{
			map[string]interface{}{},
			&map[int]int{},
			map[int]int{},
			0, true,
		},
		// Test 2
		{
			map[string]interface{}{"int_val": 1976, "float_val": 19.76},
			&struct {
				IntVal   int
				FloatVal float64
			}{},
			struct {
				IntVal   int
				FloatVal float64
			}{1976, 19.76},
			2, false,
		},
		// Test 3
		{
			map[string]interface{}{"int_val": 1976, "float_val": 19.76},
			&struct {
				intVal   int
				FloatVal float64
			}{},
			struct {
				intVal   int
				FloatVal float64
			}{0, 19.76},
			1, true,
		},
		// Test 4
		{
			map[string]interface{}{"int_val": 1976, "float_val": 19.76},
			&struct {
				Value1   int `intVal`
				FloatVal int
			}{},
			struct {
				intVal   int
				FloatVal float64
			}{1976, 0},
			1, true,
		},
	}
	for i, tt := range tests {
		n, err := ParseMapToStruct(tt.srcMap, tt.dstPtr)
		dstStruct := strings.TrimLeft(fmt.Sprint(tt.dstPtr), "&")
		exStruct := fmt.Sprint(tt.exStruct)
		if (n != tt.exNum) || ((err != nil) != tt.exErr) || (dstStruct != exStruct) {
			errStatus := "nil"
			if err != nil {
				errStatus = "error: " + err.Error()
			}
			exErrStatus := "nil"
			if tt.exErr {
				exErrStatus = "error"
			}
			t.Errorf("Test %v: ParseMapToStruct(%v, dstPtr) returned (%v, %v); expected: (%v, %v)\n"+
				"Target struct is: %v; expected: %v",
				i, tt.srcMap, n, errStatus, tt.exNum, exErrStatus, dstStruct, exStruct)
		}
	}
}
