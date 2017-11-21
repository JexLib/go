package main

import (
	"github.com/JexLib/golang/configor"
)

/**
Enabled bool   `flag:"|true|run ServiceOneConfig"`
IP      string `flag:"H|127.0.0.1|Listen {IP}" env:"PATH"`
Port    int    `flag:"p|8080|Listen {Port}"`
*/
type Config struct {
	APPName string `default:"configor"`
	Hosts   []string

	DB struct {
		Name     string
		User     string `default:"root" flag:"u|this is user"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"  flag:"p|this is db port"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}

	private string
}

func main() {
	result := Config{}
	configor.Load(&result)
	configor.PrintJson(&result)
}
