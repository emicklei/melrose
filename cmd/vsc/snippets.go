package main

// script to generate the snippets.json for melrose-for-vscode

import (
	"encoding/json"
	"os"

	"github.com/emicklei/melrose/dsl"
)

// see Makefile how to run this

type Snippet struct {
	Prefix      string   `json:"prefix"`
	Body        []string `json:"body"`
	Description string   `json:"description"`
}

func snippets() {
	varstore := dsl.NewVariableStore()
	snippets := map[string]Snippet{}
	for _, v := range dsl.EvalFunctions(varstore) {
		if len(v.Prefix) > 0 && len(v.Title) > 0 {
			snip := Snippet{
				Prefix:      v.Prefix,
				Body:        []string{v.Sample},
				Description: v.Description,
			}
			snippets[v.Title] = snip
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(snippets)
}
