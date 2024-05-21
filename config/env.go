package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

var EnvMode bool

func InitEnv() {
	if EnvMode {
		fmt.Println("Running in with env..")
		err := godotenv.Load()
		if err != nil {
			panic(err.Error())
		}
	}
}
