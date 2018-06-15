package yagolib

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

// LoadConfig loads and parses configuration file(s) in TOML format:
// https://github.com/toml-lang/toml
// Command line arguments may be parsed too if necessary.
// Configuration will be copied from data sources to the structure pointed by 'config' in the following order:
// 1. From 'homeConfigName' file in the '~/.config/appName' directory (appName - name of application executable).
// 2. From 'configPath' file.
// 3. From command line arguments given by 'cmdLine' string array.
// So 'homeConfigName' file has the lowest priority. It's settings will be overridden by 'configPath' file.
// Command line arguments have top priority and will override data from 'configPath' file.
// You may omit any config data source, just use empty string for 'homeConfigName'/'configPath' and 'nil' for 'cmdLine'.
// The fields of 'config' structure must be exported.
// The names of command line flags must exactly match the names of 'config' structure fields.
// The fields of 'config' structure may have optional `tag` which define description of flag.
// For the following structure:
// type configStruct struct {
//	  Param1   int `Parameter number One`
// }
// You can pass to command line the following flag: -Param1=35 (not -param1=35).
// The function returns '*flag.FlagSet' if 'cmdLine' was set and 'error'.
func LoadConfig(config interface{}, homeConfigName string, configPath string, cmdLine []string, verbose bool) (*flag.FlagSet, error) {
	if config == nil {
		return nil, fmt.Errorf("'config' structure pointer is 'nil'")
	}
	configType := reflect.TypeOf(config)
	if (configType.Kind() != reflect.Ptr) || (configType.Elem().Kind() != reflect.Struct) {
		return nil, fmt.Errorf("'config' argument is not a pointer to structure. It has type: %v", configType)
	}

	var errMsg string

	appName := filepath.Base(os.Args[0])

	if homeConfigName != "" {
		homeConfigPath := filepath.Join("~/.config", appName, homeConfigName)
		homeConfigPath, err := NormalizePath(homeConfigPath)
		if (err == nil) && IsFileExists(homeConfigPath) {
			if verbose {
				fmt.Printf("Loading config from '%v'\n", homeConfigPath)
			}
			if _, err = toml.DecodeFile(homeConfigPath, config); err != nil {
				errMsg = fmt.Sprintf("Error parsing config file '%v':\n%v\n", homeConfigPath, err)
				if verbose {
					fmt.Fprint(os.Stderr, errMsg)
				}
			}
		} else if verbose {
			fmt.Fprintf(os.Stderr, "Config file '%v' not found\n", homeConfigPath)
		}
	}

	if configPath != "" {
		configPath, err := NormalizePath(configPath)
		if (err == nil) && IsFileExists(configPath) {
			if verbose {
				fmt.Printf("Loading config from '%v'\n", configPath)
			}
			if _, err = toml.DecodeFile(configPath, config); err != nil {
				msg := fmt.Sprintf("Error parsing config file '%v':\n%v\n", configPath, err)
				errMsg += msg
				if verbose {
					fmt.Fprint(os.Stderr, msg)
				}
			}
		} else if verbose {
			fmt.Fprintf(os.Stderr, "Config file '%v' not found\n", configPath)
		}
	}

	var flagSet *flag.FlagSet
	if cmdLine != nil {
		var tomlStr string
		flagSet = flag.NewFlagSet(appName, flag.ContinueOnError)
		structValue := reflect.ValueOf(config).Elem()
		structType := structValue.Type()
		for i := 0; i < structType.NumField(); i++ {
			fieldValue := structValue.Field(i)
			field := structType.Field(i)
			fieldName := field.Name
			fieldTag := string(field.Tag)
			switch fieldValue.Kind() {
			case reflect.Bool:
				flagSet.Bool(fieldName, fieldValue.Bool(), fieldTag)
			case reflect.Float32:
				flagSet.Float64(fieldName, fieldValue.Float(), fieldTag)
			case reflect.Float64:
				flagSet.Float64(fieldName, fieldValue.Float(), fieldTag)
			case reflect.Int:
				flagSet.Int(fieldName, int(fieldValue.Int()), fieldTag)
			case reflect.Int64:
				flagSet.Int64(fieldName, fieldValue.Int(), fieldTag)
			case reflect.String:
				flagSet.String(fieldName, fieldValue.String(), fieldTag)
			}
		}
		if err := flagSet.Parse(cmdLine); err != nil {
			msg := fmt.Sprintf("Error parsing command line:\n%v\n", err)
			errMsg += msg
			if verbose {
				fmt.Fprint(os.Stderr, msg)
			}
		}
		flagSet.Visit(func(f *flag.Flag) {
			keyVal := fmt.Sprintf("%v = %v\n", f.Name, f.Value)
			tomlStr += keyVal
		})
		if _, err := toml.Decode(tomlStr, config); err != nil {
			msg := fmt.Sprintf("Error parsing command line:\n%v\n", err)
			errMsg += msg
			if verbose {
				fmt.Fprint(os.Stderr, msg)
			}
		}
	}

	if errMsg == "" {
		return flagSet, nil
	}
	return flagSet, fmt.Errorf(strings.TrimRight(errMsg, "\n"))
}
