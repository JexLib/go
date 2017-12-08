package main

import (
	"github.com/JexLib/golang/configor"
)

/**
Enabled bool   `flag:"|true|run ServiceOneConfig"`
IP      string `flag:"H|127.0.0.1|Listen {IP}" env:"PATH"`
Port    int    `flag:"p|8080|Listen {Port}"`

shortName|default|des|env,required
`config:"p|8080|Listen Port|PATH,required"`
*/
type Config struct {
	APPName string `default:"configor"`
	Hosts   []string

	DB struct {
		Name     string
		User     string `config:"u|root|this is user|"`
		Password string ` env:"DBPassword"`
		Port     uint   `config:"p|3306|this is db port"`
		Mm       string `config:"|m-pp|Listen Port|,required"`
		Md       string `config:"|8080|Listen Port|PATH,required"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}

	private string
}

func main() {
	result := Config{}
	//configor.Default(&result)

	configor.Load("confTest", "1.0", &result, "configB.json")
	configor.PrintJson(&result)

	// flag.Usage()
}
