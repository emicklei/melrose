package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	dir := filepath.Join(os.Getenv("HOME"), "sounds")
	list, _ := ioutil.ReadDir(dir)
	for _, each := range list {
		out := each.Name()[0 : len(each.Name())-4]
		//fmt.Printf("ffmpeg -i %s/%s -ar 8000 -ac 1 %s/%s\n", dir, each.Name(), dir, out)
		fmt.Printf("lame --decode %s/%s  %s/%s\n", dir, each.Name(), dir, out)
	}
}