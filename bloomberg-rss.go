package main

import (
	"github.com/posativ/bloomberg-rss/src/server"
)

func main() {
	app, err := server.NewServer()
	if err != nil {
		panic(err)
	}
	app.Start()
}
