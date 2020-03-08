package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Check the license of multiplayerpiano.com before fetching their WAV-es.

func main() {
	dir := filepath.Join(os.Getenv("HOME"), "sounds")
	fmt.Println("writing sounds in", dir)
	base := "http://multiplayerpiano.com/sounds/mppclassic/"
	notes := []string{"c", "cs", "d", "ds", "e", "f", "fs", "g", "gs", "a", "as", "b"}
	for oct := 0; oct < 7; oct++ {
		for _, note := range notes {
			prefix := note + fmt.Sprintf("%d", oct)
			full := base + prefix + ".mp3"
			resp, _ := http.Get(full)
			fmt.Printf("fetching:%s\n", full)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			ioutil.WriteFile(filepath.Join(dir, prefix+".mp3"), body, os.ModePerm)
		}
	}
}
