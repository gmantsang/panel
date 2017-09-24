package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/dabbotorg/panel/config"
	"github.com/dabbotorg/panel/server"
)

var conf config.Config
var meta config.Metadata

func init() {
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = toml.DecodeFile("metadata.toml", &meta)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println("panel up")
	log.Println(server.Serve(conf, meta))
	log.Println("panel down")
}
