package smartConfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	// 默认的 struct tag key，本模块据此从 struct tag 中寻找需要的配置
	// 例如:
	// type Config struct {
	//     Foo string    `flag:"...Foo 字段的命令行选项属性..."`
	// }
	StructTagKey_flag string = "flag"
	StructTagKey_env  string = "env"
)

var (
	configFile    string
	configFileOK  bool
	configMonitor chan interface{}
)

func LoadConfig(name string, version string, config interface{}) {
	cobra.OnInitialize(initConfig)

	configMonitor = make(chan interface{}, 1)

	cmdRoot := &cobra.Command{
		Use:   os.Args[0],
		Short: name,
		Long:  fmt.Sprintf("%s(version %s)", name, version),
		Run: func(*cobra.Command, []string) {
			loadConfig(config)
			if configFileOK {
				viper.WatchConfig()
				viper.OnConfigChange(func(e fsnotify.Event) {
					loadConfig(config)
					noticeChanged(config)
				})
			}
		},
	}

	flags := cmdRoot.PersistentFlags()
	cmdRoot.Flags().SortFlags = false
	cmdRoot.PersistentFlags().SortFlags = false

	flags.StringVarP(&configFile, "config", "c", "", "config `FILENAME`, default to `config.yaml` or `config.json`")
	optVersion := flags.Bool("version", false, "just print version number only")
	optHelp := flags.Bool("help", false, "show this message")

	optWriteYAML := flags.Bool("gen-yaml", false, "generate config.example.yaml")
	optWriteJSON := flags.Bool("gen-json", false, "generate config.example.json")

	addFlags(flags, config)

	viper.BindPFlags(flags)

	err := cmdRoot.Execute()
	if err != nil || *optHelp {
		os.Exit(0)
	} else if *optVersion {
		fmt.Fprintf(os.Stderr, "%s\n", version)
		os.Exit(0)
	} else if *optWriteYAML {
		// out, err := yaml.Marshal(config)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "error: %s\n", err)
		// } else {
		// 	//fmt.Fprintf(os.Stdout, "%s", out)
		// 	createExampleConfigFile(config, "yaml")
		// }
		gen_ConfigFile(config, "yaml")
		os.Exit(0)
	} else if *optWriteJSON {
		// encoder := json.NewEncoder(os.Stdout)
		// encoder.SetIndent("", "    ")
		// encoder.Encode(config)
		gen_ConfigFile(config, "json")
		os.Exit(0)
	}
}

type FlagSpec struct {
	Type        reflect.Type
	Name        string
	ShortName   string
	Default     string
	HelpMessage string
}

func addFlags(flags *pflag.FlagSet, config interface{}) {
	list := getList("", reflect.TypeOf(config), "", "")
	used := make(map[string]string, 32)
	for _, v := range list {
		help := v.HelpMessage
		begin := strings.Index(help, "{")
		end := strings.Index(help, "}")
		if begin >= 0 && end > begin {
			help = help[0:begin] + "`" + help[begin+1:end] + "`" + help[end+1:]
		}

		switch len(v.ShortName) {
		case 0:
		case 1:
			if used[v.ShortName] != "" {
				fmt.Fprintf(os.Stderr, "选项 %s 的短名称 %s 不能生效，因为已经被 %s 占用。\n", v.Name, v.ShortName, used[v.ShortName])
				v.ShortName = ""
			} else {
				used[v.ShortName] = v.Name
			}
		default:
			fmt.Fprintf(os.Stderr, "选项 %s 的短名称 %s 不能生效，因为长度超过了一个字节。\n", v.Name, v.ShortName)
			v.ShortName = ""
		}

		if v.Type == reflect.TypeOf(time.Second) {
			value, err := time.ParseDuration(v.Default)
			if err != nil {
				fmt.Fprintf(os.Stderr, "选项 %s 的默认值 %s 的格式不对。\n", v.Name, v.Default)
			}
			flags.DurationP(v.Name, v.ShortName, value, help)
			continue
		}

		switch v.Type.Kind() {
		case reflect.Bool:
			var value bool
			fmt.Sscanf(v.Default, "%v", &value)
			flags.BoolP(v.Name, v.ShortName, value, help)
		case reflect.Int:
			var value int
			fmt.Sscanf(v.Default, "%v", &value)
			flags.IntP(v.Name, v.ShortName, value, help)
		case reflect.Int8:
			var value int8
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Int8P(v.Name, v.ShortName, value, help)
		case reflect.Int16:
			var value int32
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Int32P(v.Name, v.ShortName, value, help)
		case reflect.Int32:
			var value int32
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Int32P(v.Name, v.ShortName, value, help)
		case reflect.Int64:
			var value int64
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Int64P(v.Name, v.ShortName, value, help)
		case reflect.Uint:
			var value uint
			fmt.Sscanf(v.Default, "%v", &value)
			flags.UintP(v.Name, v.ShortName, value, help)
		case reflect.Uint8:
			var value uint8
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Uint8P(v.Name, v.ShortName, value, help)
		case reflect.Uint16:
			var value uint16
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Uint16P(v.Name, v.ShortName, value, help)
		case reflect.Uint32:
			var value uint32
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Uint32P(v.Name, v.ShortName, value, help)
		case reflect.Uint64:
			var value uint64
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Uint64P(v.Name, v.ShortName, value, help)
		case reflect.Float32:
			var value float32
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Float32P(v.Name, v.ShortName, value, help)
		case reflect.Float64:
			var value float64
			fmt.Sscanf(v.Default, "%v", &value)
			flags.Float64P(v.Name, v.ShortName, value, help)
		case reflect.String:
			var value string
			fmt.Sscanf(v.Default, "%v", &value)
			flags.StringP(v.Name, v.ShortName, value, help)
		}
	}
}

func getList(prefix string, t reflect.Type, tag_flag, tag_env string) (result []FlagSpec) {
	result = make([]FlagSpec, 0)

	if t == reflect.TypeOf(time.Second) {
		goto end
	}

	switch t.Kind() {
	case reflect.Ptr:
		result = getList(prefix, t.Elem(), tag_flag, tag_env)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i += 1 {
			tag_flag := t.Field(i).Tag.Get(StructTagKey_flag)
			tag_env := t.Field(i).Tag.Get(StructTagKey_env)
			name := prefix
			if name != "" {
				name += "."
			}
			name += strings.ToLower(t.Field(i).Name)
			result = append(result, getList(name, t.Field(i).Type, tag_flag, tag_env)...)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Bool, reflect.String, reflect.Float32, reflect.Float64:
		goto end
	default:
	}

end:
	if tag_flag == "" && tag_env == "" {
		return
	}

	parts := strings.SplitN(tag_flag, "|", 3)
	if tag_env != "" {
		if value := os.Getenv(tag_env); value != "" {
			parts[1] = value
		}
	}

	result = append(result, FlagSpec{
		Type:        t,
		Name:        prefix,
		ShortName:   parts[0],
		Default:     parts[1],
		HelpMessage: parts[2],
	})

	return
}

func initConfig() {
	if configFile != "" {
		// 允许通过命令行参数来指定配置文件路径
		viper.SetConfigFile(configFile)
	} else {
		// 否则在当前目录下寻找 config.{yaml,yml,json} 等文件
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	err := viper.ReadInConfig()
	if err == nil {
		configFileOK = true
	}
}

func loadConfig(config interface{}) error {
	err := viper.Unmarshal(config)
	if err != nil {
		log.Print("unable to decode into struct: ", err)
		return err
	}

	return nil
}

func noticeChanged(config interface{}) {
	select {
	case configMonitor <- config:
	default:
	}
}

func ConfigChanged() <-chan interface{} {
	return configMonitor
}

/**
 configType:yaml,json
**/
func gen_ConfigFile(config interface{}, configType ...string) {
	if len(configType) == 0 {
		configType = append(configType, "json")
	}
	switch configType[0] {
	case "yaml":
		out, err := yaml.Marshal(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}
		ioutil.WriteFile("./config.example.yaml", out, 0666)
		fmt.Println("generate ./config.example.yaml succeed:\n", string(out))
	case "json":
		out, err := json.Marshal(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		var fmtout bytes.Buffer
		json.Indent(&fmtout, out, "", "\t")
		ioutil.WriteFile("./config.example.json", fmtout.Bytes(), 0666)
		fmt.Println("generate ./config.example.json succeed:\n", fmtout.String())
	}

}
