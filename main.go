package main

import (
	"caminoclient/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
