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
	base := "http://multiplayerpiano.com/mp3/"
	notes := []string{"c", "cs", "d", "ds", "e", "f", "fs", "g", "gs", "a", "as", "b"}
	for oct := 1; oct < 7; oct++ {
		for _, note := range notes {
			prefix := note + fmt.Sprintf("%d", oct)
			resp, _ := http.Get(base + prefix + ".wav.mp3")
			fmt.Printf("fetching:%s\n", prefix)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			ioutil.WriteFile(filepath.Join(dir, prefix+".wav.mp3"), body, os.ModePerm)
		}
	}
}
