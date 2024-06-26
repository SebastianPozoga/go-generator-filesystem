package main

import (
	"flag"
)

func main() {
	var app = &App{}
	flag.StringVar(&app.From, "from", "", "directory or file to convert from")
	flag.StringVar(&app.To, "to", "", "destination directory or file")
	flag.StringVar(&app.Cache, "cache", "", "destination directory for cache data")
	flag.StringVar(&app.PackagePrefix, "package.prefix", "", "prefix for all package imports")
	flag.BoolVar(&app.LogAll, "logs", false, "view full log")
	flag.Parse()

	app.Valid()
	app.InitFS()
	app.Run()
}
