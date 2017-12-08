package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/JexLib/golang/crypto/base58"
)

func main() {
	if len(os.Args) < 3 {
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "e", "-e":
		fmt.Println("Encode:", base58.Encode([]byte(os.Args[2])))
	case "d", "-d":
		fmt.Println("Decode:", string(base58.Decode(os.Args[2])))
	}

}
