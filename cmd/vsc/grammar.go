package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/emicklei/melrose/core"

	"github.com/emicklei/melrose/dsl"
)

// see Makefile how to run this

func grammar() {
	data, err := os.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	content := string(data)

	varstore := dsl.NewVariableStore()
	ctx := core.PlayContext{
		VariableStorage: varstore,
		LoopControl:     core.NoLooper,
	}
	// collect and sort by length , descending
	keywords := []string{}
	for k, f := range dsl.EvalFunctions(ctx) {
		keywords = append(keywords, k)
		if len(f.Alias) > 0 {
			keywords = append(keywords, f.Alias)
		}
	}
	sort.SliceStable(keywords, func(i, j int) bool {
		return len(keywords[i]) > len(keywords[j])
	})
	var buffer bytes.Buffer
	for _, k := range keywords {
		if buffer.Len() > 0 {
			fmt.Fprintf(&buffer, "|")
		}
		fmt.Fprintf(&buffer, "%s", k)
	}
	content = strings.Replace(content, "$Keywords", buffer.String(), -1)
	if err := os.WriteFile(os.Args[3], []byte(content), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
