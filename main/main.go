package main

import (
	"github.com/sid-sun/arche-api/app"
	"github.com/sid-sun/arche-api/config"
	"github.com/sid-sun/arche-api/ginit"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	lgr, err := ginit.Logger(cfg)
	if err != nil {
		panic(err)
	}

	app.Start(cfg, lgr)
}
