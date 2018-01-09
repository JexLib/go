package configor

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

func (configor *Configor) getENVPrefix(config interface{}) string {
	if configor.Config.ENVPrefix == "" {
		if prefix := os.Getenv("CONFIGOR_ENV_PREFIX"); prefix != "" {
			return prefix
		}
		return "Configor"
	}
	return configor.Config.ENVPrefix
}

func getConfigurationFileWithENVPrefix(file, env string) (string, error) {
	var (
		envFile string
		extname = path.Ext(file)
	)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
		return envFile, nil
	}
	return "", fmt.Errorf("failed to find file %v", file)
}

func (configor *Configor) getConfigurationFiles(files ...string) []string {
	var results []string

	if configor.Config.Debug || configor.Config.Verbose {
		fmt.Printf("Current environment: '%v'\n", configor.GetEnvironment())
	}

	for i := len(files) - 1; i >= 0; i-- {
		foundFile := false
		file := files[i]

		// check configuration
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			foundFile = true
			results = append(results, file)
		}

		// check configuration with env
		if file, err := getConfigurationFileWithENVPrefix(file, configor.GetEnvironment()); err == nil {
			foundFile = true
			results = append(results, file)
		}

		// check example configuration
		if !foundFile {
			if example, err := getConfigurationFileWithENVPrefix(file, "example"); err == nil {
				fmt.Printf("Failed to find configuration %v, using example file %v\n", file, example)
				results = append(results, example)
			} else {
				fmt.Printf("Failed to find configuration %v\n", file)
			}
		}
	}
	return results
}

func processFile(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		return yaml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".toml"):
		return toml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".json"):
		return json.Unmarshal(data, config)
	default:
		if toml.Unmarshal(data, config) != nil {
			if json.Unmarshal(data, config) != nil {
				if yaml.Unmarshal(data, config) != nil {
					return errors.New("failed to decode config")
				}
			}
		}
		return nil
	}
}

func getPrefixForStruct(prefixes []string, fieldStruct *reflect.StructField) []string {
	if fieldStruct.Anonymous && fieldStruct.Tag.Get("anonymous") == "true" {
		return prefixes
	}
	return append(prefixes, fieldStruct.Name)
}

/**
shortName|default|HelpMessage|env,required
`config:"p|8080|Listen Port|PATH,required"`
*/
func (configor *Configor) processTags(createFlag bool, config interface{}, ntag []string, prefixes ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			envNames         []string
			fieldStruct      = configType.Field(i)
			field            = configValue.Field(i)
			fieldTag         = fieldStruct.Tag.Get("config")
			fieldRequired    = false
			fieldShortName   = ""
			fieldDefault     = ""
			fieldHelpMessage = ""
			fieldEnvName     = "" // read configuration from shell env
		)

		for _, v := range strings.Split(fieldTag, ",") {
			switch v {
			case "required":
				fieldRequired = true
			default:
				tags := strings.Split(v, "|")
				for len(tags) < 4 {
					tags = append(tags, "")
				}
				fieldShortName = tags[0]
				fieldDefault = tags[1]
				fieldHelpMessage = tags[2]
				fieldEnvName = tags[3]
			}

		}

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		if fieldEnvName == "" {
			envNames = append(envNames, strings.Join(append(prefixes, fieldStruct.Name), "_"))                  // Configor_DB_Name
			envNames = append(envNames, strings.ToUpper(strings.Join(append(prefixes, fieldStruct.Name), "_"))) // CONFIGOR_DB_NAME
		} else {
			envNames = []string{fieldEnvName}
		}

		if configor.Config.Verbose {
			fmt.Printf("Trying to load struct `%v`'s field `%v` from env %v\n", configType.Name(), fieldStruct.Name, strings.Join(envNames, ", "))
		}

		// Load From Shell ENV
		for _, env := range envNames {
			if value := os.Getenv(env); value != "" {
				if configor.Config.Debug || configor.Config.Verbose {
					fmt.Printf("Loading configuration for struct `%v`'s field `%v` from env %v...\n", configType.Name(), fieldStruct.Name, env)
				}
				if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
					return err
				}
				break
			}
		}

		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); !createFlag && isBlank {
			// Set default configuration if blank
			if fieldDefault != "" {
				if err := yaml.Unmarshal([]byte(fieldDefault), field.Addr().Interface()); err != nil {
					return err
				}
			} else if fieldRequired {
				// return error if it is required but blank
				return errors.New(strings.Join(append(ntag, fieldStruct.Name), ".") + " is required, but blank")
			}
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {

			if err := configor.processTags(createFlag, field.Addr().Interface(), append(ntag, fieldStruct.Name), getPrefixForStruct(prefixes, &fieldStruct)...); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Slice {
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := configor.processTags(createFlag, field.Index(i).Addr().Interface(), append(ntag, fieldStruct.Name), append(getPrefixForStruct(prefixes, &fieldStruct), fmt.Sprint(i))...); err != nil {
						return err
					}
				}
			}
		}

		if createFlag && field.Kind() != reflect.Struct {
			configor.flagValue(flagSpec{
				Type:        field.Type(),
				Name:        strings.Join(append(ntag, fieldStruct.Name), "."),
				ShortName:   fieldShortName,
				Default:     fieldDefault,
				HelpMessage: fieldHelpMessage,
				Required:    fieldRequired,
				Env:         fieldEnvName,
				Ptr:         field.Addr().Interface(),
			})
		}
	}
	return nil
}

// func usage(name string, version string) {
// 	fmt.Fprintf(os.Stderr, name+` version: `+version+`
// Usage: `+os.Args[0]+` [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

// Options:
// `)
// 	flag.PrintDefaults()
// }

func (configor *Configor) flagUsage() {

	fmt.Println("Usage:")
	fmt.Println(" ", os.Args[0], "[OPTIONS]\n")
	fmt.Println("Application Options:")
	fmt.Println(" ", "-C, --configfile=       Path to configuration file (/root/.stakepoold/stakepoold.conf")
	fmt.Println(" ", "-V, --version           Display version information and exit")
	fmt.Println()
	fmt.Println("Help Options:")
	fmt.Println(" ", "-h, --help              Show this help message")
	f := flag.CommandLine

	f.VisitAll(func(flag *flag.Flag) {
		s := fmt.Sprintf("  -%s", flag) // Two spaces before -; see next two comments.
		fmt.Println(s)
		// name, usage := flag.UnquoteUsage(flag)
		// if len(name) > 0 {
		// 	s += " " + name
		// }
		// // Boolean flags of one ASCII letter are so common we
		// // treat them specially, putting their usage on the same line.
		// if len(s) <= 4 { // space, space, '-', 'x'.
		// 	s += "\t"
		// } else {
		// 	// Four spaces before the tab triggers good alignment
		// 	// for both 4- and 8-space tab stops.
		// 	s += "\n    \t"
		// }
		// s += usage
		// if !flag.isZeroValue(flag, flag.DefValue) {
		// 	if _, ok := flag.Value.(*stringValue); ok {
		// 		// put quotes on the value
		// 		s += fmt.Sprintf(" (default %q)", flag.DefValue)
		// 	} else {
		// 		s += fmt.Sprintf(" (default %v)", flag.DefValue)
		// 	}
		// }
		// fmt.Fprint(f.out(), s, "\n")
	})
}

func PrintJson(config interface{}) {
	var fmtout bytes.Buffer
	out, err := json.Marshal(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	json.Indent(&fmtout, out, "", "\t")
	fmt.Println(fmtout.String())
}
