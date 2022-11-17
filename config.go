package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type confItem struct {
	Res         *string
	Default     string
	PanicIfNone bool
}

var (
	DB_URI  = ""
	PORT    = ""
)

func InitConfig() {
	godotenv.Load(".env")

	var confMap = map[string]confItem{
		"DB_URI": {
			Res:         &DB_URI,
			PanicIfNone: true,
		},
		"PORT": {
			Res:         &PORT,
			Default: "3000",
		},
	}

	for name, opt := range confMap {
		item := os.Getenv(name)

		if item == "" {
			if opt.PanicIfNone {
				panic(fmt.Sprintf("'%v' is a needed variable, but is not present! Please read the README.md file for more info.", name))
			}
			item = opt.Default
		}

		*opt.Res = item
	}
}
