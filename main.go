package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Constants
var REGISTRY string = "https://registry.npmjs.org/"

// var CACHE string = "./.cache"
var MODULES string = "./node_modules"
var MODULESBIN string = filepath.Join(MODULES, ".bin")

// var workDir string
var InstalledPackages int = 0
var FetchedPackages int = 0

func main() {
	_, err := os.Stat(MODULES)
	if os.IsNotExist(err) {
		if err := os.Mkdir(MODULES, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// _, er := os.Stat(CACHE)
	// if os.IsNotExist(er) {
	// 	if err := os.Mkdir(CACHE, os.ModePerm); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// then := time.Now()
	// workDir = os.Args[0]
	command := os.Args[1:]
	resolveCmd(strings.Join(command, " "))
}
