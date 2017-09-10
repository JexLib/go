package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	smartConfig "github.com/JexLib/go/go-smart-config"
)

func main() {
	// 实际应用中建议将 ServiceOneConfig 等类型的定义放到具体的业务模块当中
	type ServiceOneConfig struct {
		Enabled bool   `flag:"|true|run ServiceOneConfig"`
		IP      string `flag:"H|127.0.0.1|Listen {IP}" env:"PATH"`
		Port    int    `flag:"p|8080|Listen {Port}"`
	}

	// 实际应用中建议将 ServiceTwoConfig 等类型的定义放到具体的业务模块当中
	type ServiceTwoConfig struct {
		Foo bool   `flag:"|true|help message for foo"`
		Bar string `flag:"|blablabla|help message for bar"`
	}

	// 实际应用中建议将 ServiceThreeConfig 等类型的定义放到具体的业务模块当中
	type ServiceThreeConfig struct {
		Hello int32         `flag:"|100|help message for hello"`
		World time.Duration `flag:"|30s|help message for world"`
	}

	type MainConfig struct {
		Debug bool `flag:"v|false|debug mode"`
		One   ServiceOneConfig
		Two   ServiceTwoConfig
		Three ServiceThreeConfig
	}

	Config := MainConfig{}
	//创建配置文件样本
	//smartConfig.CreateExampleConfigFile(&Config)
	//读取配置
	smartConfig.LoadConfig("example.yml", "1.0", &Config)

	if str, err := json.Marshal(&Config); err == nil {
		var out bytes.Buffer
		json.Indent(&out, str, "", "\t")
		fmt.Println("sss:", out.String())
	}

	for {
		select {
		case <-smartConfig.ConfigChanged():
			// 监测到配置文件修改
			fmt.Printf("new config: %#v\n", Config)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
