package main

import (
	"fmt"

	"github.com/JexLib/golang/configor"
)

/**
Enabled bool   `flag:"|true|run ServiceOneConfig"`
IP      string `flag:"H|127.0.0.1|Listen {IP}" env:"PATH"`
Port    int    `flag:"p|8080|Listen {Port}"`

shortName|default|des|env,required
`config:"p|8080|Listen Port|PATH,required"`
*/
type Sub struct {
	SS string `config:"|sss|sssssssssss"`
}

type Config struct {
	APPName string `default:"configor"`
	Hosts   []string

	Sub Sub `config:" required"`
	DB  struct {
		Name     string
		User     string `config:"u|root|this is user|"`
		Password string `config:"|8888|DBPassword"`
		Port     uint   `config:"p|3306|this is db port"`
		Mm       string `config:"|m-pp|Listen Port|,required"`
		Md       string `config:"|8080|Listen Port|PATH,required"`
	}

	Contacts []struct {
		Name  string
		Email string
	}

	private string
}

func main() {
	result := Config{}
	configor.Default(&result)
	fmt.Println("default:")
	configor.PrintJson(&result)

	if err := configor.Load("confTest", "1.0", &result, "configB.json"); err == nil {
		fmt.Println("loads:")
		configor.PrintJson(&result)
	} else {
		fmt.Println("loads err:", err)
	}

	rr := Config{
		APPName: "name11",
	}
	configor.Default(&rr)
	configor.PrintJson(&rr)
	// flag.Usage()
}
