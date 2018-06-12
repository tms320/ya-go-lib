package yagolib

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/BurntSushi/toml"
)

// LoadConfig loads and parses configuration file(s) in TOML format:
// https://github.com/toml-lang/toml
// Command line arguments are parsed too if necessary.
// Configuration data will be copied to 'config' structure in the following order:
// 1. From 'homeConfigName' file in the '~/.config/appName' directory (appName - name of application executable).
// 2. 'configPath' file.
// 3. Command line arguments.
// So 'homeConfigName' file has the lowest priority. It's settings can be overridden by 'configPath' file.
// Command line arguments have top priority and will override data from 'configPath' file.
// You may omit any config data source (just use empty string for homeConfigName/configPath).
// If you want to parse command line flags, the fields of target 'config' structure may
// have struct `tags` which define alternative name and description of flag like this:
// type configStruct struct {
//	  Param1   int `param_1: Parameter number One`
// }
// So you can pass to command line the following variants of flag name:
// -Param1=35 or -param1=35 or -param_1=35
func LoadConfig(config interface{}, homeConfigName string, configPath string, parseCmdLine bool, verbose bool) {
	if config == nil {
		fmt.Fprintln(os.Stderr, "Error: target 'config' structure is 'nil'")
		return
	}
	configType := reflect.TypeOf(config)
	if configType.Kind() != reflect.Ptr {
		fmt.Fprintln(os.Stderr, "Error: target 'config' structure is not a pointer. It has type: ", configType)
		return
	}

	if homeConfigName != "" {
		exePath, err := os.Executable()
		if err == nil {
			appName := filepath.Base(exePath)
			homeConfigPath := filepath.Join("~/.config", appName, homeConfigName)
			homeConfigPath, err = NormalizePath(homeConfigPath)
			if (err == nil) && IsFileExists(homeConfigPath) {
				if verbose {
					fmt.Printf("Loading config from '%v'\n", homeConfigPath)
				}
				if _, err = toml.DecodeFile(homeConfigPath, config); (err != nil) && verbose {
					fmt.Fprintf(os.Stderr, "Error parsing config file '%v':\n", homeConfigPath)
					fmt.Fprint(os.Stderr, err)
				}
			} else if verbose {
				fmt.Fprintf(os.Stderr, "Config file '%v' not found\n", homeConfigPath)
			}
		}
	}

	if configPath != "" {
		configPath, err := NormalizePath(configPath)
		if (err == nil) && IsFileExists(configPath) {
			if verbose {
				fmt.Printf("Loading config from '%v'\n", configPath)
			}
			if _, err = toml.DecodeFile(configPath, config); (err != nil) && verbose {
				fmt.Fprintf(os.Stderr, "Error parsing config file '%v':\n", configPath)
				fmt.Fprint(os.Stderr, err)
			}
		} else if verbose {
			fmt.Fprintf(os.Stderr, "Config file '%v' not found\n", configPath)
		}
	}

	if parseCmdLine {
		flag.Parse()
	}

}
