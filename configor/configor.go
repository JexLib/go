package configor

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Configor struct {
	*Config
}

type Config struct {
	Environment string
	ENVPrefix   string
	Debug       bool
	Verbose     bool
}

type flagSpec struct {
	Type        reflect.Type
	Name        string
	ShortName   string
	Default     string
	HelpMessage string
	Required    bool
	Env         string
	Ptr         interface{}
}

// New initialize a Configor
func New(config *Config) *Configor {
	if config == nil {
		config = &Config{}
	}

	if os.Getenv("CONFIGOR_DEBUG_MODE") != "" {
		config.Debug = true
	}

	if os.Getenv("CONFIGOR_VERBOSE_MODE") != "" {
		config.Verbose = true
	}

	return &Configor{Config: config}
}

// GetEnvironment get environment
func (configor *Configor) GetEnvironment() string {
	if configor.Environment == "" {
		if env := os.Getenv("CONFIGOR_ENV"); env != "" {
			return env
		}

		if isTest, _ := regexp.MatchString("/_test/", os.Args[0]); isTest {
			return "test"
		}

		return "development"
	}
	return configor.Environment
}

func (configor *Configor) flagValue(v flagSpec) {
	fmt.Println(v.Name)
	help := v.HelpMessage
	if v.Required {
		help = "<required> " + v.HelpMessage
	}

	switch v.Type.Kind() {
	case reflect.Bool:
		var value bool
		fmt.Sscanf(v.Default, "%v", &value)
		v.Ptr = flag.Bool(v.ShortName, value, help)
	case reflect.Int:
		var value int
		fmt.Sscanf(v.Default, "%v", &value)
		v.Ptr = flag.Int(v.ShortName, value, help)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var value int64
		fmt.Sscanf(v.Default, "%v", &value)
		v.Ptr = flag.Int64(v.ShortName, value, help)
	case reflect.Uint:
		var value uint
		fmt.Sscanf(v.Default, "%v", &value)
		v.Ptr = flag.Uint(v.ShortName, value, help)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var value uint64
		fmt.Sscanf(v.Default, "%v", &value)
		v.Ptr = flag.Uint64(v.ShortName, value, help)
	case reflect.Float32, reflect.Float64:
		var value float64
		fmt.Sscanf(v.Default, "%v", &value)
		v.Ptr = flag.Float64(v.ShortName, value, help)
	case reflect.String:
		var value string
		fmt.Sscanf(v.Default, "%v", &value)
		v.Ptr = flag.String(v.ShortName, value, help)
	}
}

func (configor *Configor) flagFromTag(config interface{}, ntag []string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}
	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}
		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		if field.Kind() == reflect.Struct {
			if err := configor.flagFromTag(field.Addr().Interface(), append(ntag, fieldStruct.Name)); err != nil {
				return err
			}
			continue
		}

		fmt.Println(strings.Join(ntag, ".") + "." + fieldStruct.Name)

		if field.Kind() == reflect.Slice {
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := configor.flagFromTag(field.Index(i).Addr().Interface(), append(ntag, fieldStruct.Name)); err != nil {
						return err
					}
				}
			}
		}

		// fmt.Println(strings.Join(ntag, "."))
		if fieldStruct.Tag.Get("flag") == "" {
			continue
		}

		configor.flagValue(flagSpec{
			Type:        field.Type(),
			Name:        strings.Join(ntag, "."),
			ShortName:   strings.Split(fieldStruct.Tag.Get("flag"), "|")[0],
			Default:     fieldStruct.Tag.Get("default"),
			HelpMessage: strings.Split(fieldStruct.Tag.Get("flag"), "|")[1],
			Required:    fieldStruct.Tag.Get("required") == "true",
			Env:         fieldStruct.Tag.Get("env"),
			Ptr:         field.Addr().Interface(),
		})
	}
	return nil
}

func (configor *Configor) flagSet(config interface{}, files []string) {
	optHelp := flag.Bool("h", false, "show this `help` message")
	optFile := flag.String("c", "conf/config.json", "set configuration `file`")
	configor.flagFromTag(config, []string{})
	flag.Parse()
	if *optHelp {
		flag.Usage()
		os.Exit(0)
	}

	if optFile != nil {
		if len(files) > 0 {
			files = strings.Split(*optFile, ",")
		}
	}
}

// Load will unmarshal configurations to struct from files that you provide
func (configor *Configor) Load(config interface{}, files ...string) (err error) {
	defer func() {

		if configor.Config.Debug || configor.Config.Verbose {
			fmt.Printf("Configuration:\n  %#v\n", config)
		}
	}()
	configor.flagSet(config, files)
	for _, file := range configor.getConfigurationFiles(files...) {
		if configor.Config.Debug || configor.Config.Verbose {
			fmt.Printf("Loading configurations from file '%v'...\n", file)
		}
		if err = processFile(config, file); err != nil {
			return err
		}
	}

	prefix := configor.getENVPrefix(config)
	if prefix == "-" {
		return configor.processTags(config)
	}
	return configor.processTags(config, prefix)
}

// ENV return environment
func ENV() string {
	return New(nil).GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return New(nil).Load(config, files...)
}
