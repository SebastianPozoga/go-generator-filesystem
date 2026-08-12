package main

import (
	"flag"
	"strings"

	"github.com/SebastianPozoga/go-generator-filesystem/fsapp"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		*s = append(*s, item)
	}
	return nil
}

func main() {
	var app = &fsapp.App{}
	flag.StringVar(&app.From, "from", "", "directory or file to convert from")
	flag.StringVar(&app.To, "to", "", "destination directory or file")
	flag.StringVar(&app.Cache, "cache", "", "destination directory for cache data")
	flag.StringVar(&app.PackagePrefix, "package.prefix", "", "prefix for all package imports")
	flag.BoolVar(&app.LogAll, "logs", false, "view full log")
	flag.Var((*stringSliceFlag)(&app.IgnoredDirs), "ignore.dirs", "comma-separated or repeatable directory names/paths to skip")
	flag.Var((*stringSliceFlag)(&app.IgnoredFiles), "ignore.files", "comma-separated or repeatable file names/paths to skip")
	flag.Parse()

	app.Valid()
	app.InitFS()
	app.Run()
}
