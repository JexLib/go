package configor

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type Configor struct {
	*Config

	optHelp *bool
	optFile *string
	optVer  *bool
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
	help := v.HelpMessage
	if v.Required {
		help = "<*required> " + v.HelpMessage
	}
	if v.ShortName == "" {
		v.ShortName = v.Name
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

//根据tag获取默认结构数据
func (configor *Configor) Default(config interface{}) (err error) {
	prefix := configor.getENVPrefix(config)
	if prefix == "-" {
		err = configor.processTags(false, config, []string{})
	}
	err = configor.processTags(false, config, []string{}, prefix)
	return
}

func (configor *Configor) _Load(name string, version string, createFlag bool, config interface{}, files ...string) (err error) {

	if !createFlag {
		// for _, file := range configor.getConfigurationFiles(files...) {
		for _, file := range files {
			if configor.Config.Debug || configor.Config.Verbose {
				fmt.Printf("Loading configurations from file '%v'...\n", file)
			}
			if err = processFile(config, file); err != nil {
				return err
			}
		}
	}

	if createFlag {
		configor.optHelp = flag.Bool("h", false, "show this `help` message")
		//"conf/config.json"
		configor.optFile = flag.String("c", strings.Join(files, ","), "set multiple configuration `file` \n            example:conf/config-main.json,conf/config-api.json")
		configor.optVer = flag.Bool("v", false, "view `version` ")
	}

	prefix := configor.getENVPrefix(config)
	if prefix == "-" {
		err = configor.processTags(createFlag, config, []string{})
	}
	err = configor.processTags(createFlag, config, []string{}, prefix)

	return err
}

// Load will unmarshal configurations to struct from files that you provide
func (configor *Configor) Load(name string, version string, config interface{}, files ...string) (err error) {
	if err = configor._Load(name, version, true, config, files...); err == nil {

		// 改变默认的 Usage
		flag.Usage = configor.flagUsage
		flag.Parse()
		if *configor.optHelp {
			///fmt.Println("hhhhhhhhhhhhhhhhhhhhhhhhhh")
			flag.Usage()

			os.Exit(0)
		}
		if *configor.optVer {
			fmt.Println(name + " version: " + version + "  " + runtime.GOOS + "/" + runtime.GOARCH)
			os.Exit(0)
		}
		//if configor.optFile != nil {
		//	if *configor.optFile != strings.Join(files, ",") {
		afiles := strings.Split(*configor.optFile, ",")
		return configor._Load(name, version, false, config, afiles...)
		//	}

		//}

		return
	}

	// defer func() {

	// 	if configor.Config.Debug || configor.Config.Verbose {
	// 		fmt.Printf("Configuration:\n  %#v\n", config)
	// 	}
	// }()
	// //	configor.flagSet(config, files)
	// for _, file := range configor.getConfigurationFiles(files...) {
	// 	if configor.Config.Debug || configor.Config.Verbose {
	// 		fmt.Printf("Loading configurations from file '%v'...\n", file)
	// 	}
	// 	if err = processFile(config, file); err != nil {
	// 		return err
	// 	}
	// }

	// prefix := configor.getENVPrefix(config)
	// if prefix == "-" {
	// 	return configor.processTags(config, []string{})
	// }
	// return configor.processTags(config, []string{}, prefix)
	return
}

// ENV return environment
func ENV() string {
	return New(nil).GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(name string, version string, config interface{}, files ...string) error {
	return New(nil).Load(name, version, config, files...)
}

//根据tag获取默认结构数据
func Default(config interface{}) (err error) {
	return New(nil).Default(config)
}
