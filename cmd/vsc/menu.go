package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func postProcessMenus() {
	data, err := os.ReadFile("../../docs/menu.template")
	checkErr(err)
	mds, err := filepath.Glob("../../docs/*.md")
	checkErr(err)
	for _, each := range mds {
		postProcessMenu(each, string(data))
	}
}

func postProcessMenu(file string, menu string) {
	data, err := os.ReadFile(file)
	checkErr(err)
	replaced := strings.Replace(string(data), "$$menu", menu, -1)
	err = os.WriteFile(file, []byte(replaced), os.ModePerm)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
