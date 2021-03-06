package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func postProcessMenus() {
	data, err := ioutil.ReadFile("../../docs/menu.template")
	checkErr(err)
	mds, err := filepath.Glob("../../docs/*.md")
	checkErr(err)
	for _, each := range mds {
		postProcessMenu(each, string(data))
	}
}

func postProcessMenu(file string, menu string) {
	data, err := ioutil.ReadFile(file)
	checkErr(err)
	replaced := strings.Replace(string(data), "$$menu", menu, -1)
	err = ioutil.WriteFile(file, []byte(replaced), os.ModePerm)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
