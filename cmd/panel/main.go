package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/dabbotorg/panel/config"
	"github.com/dabbotorg/panel/server"
)

var conf config.Config

func init() {
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println("panel up")
	log.Println(server.Serve(conf))
	log.Println("panel down")
}
