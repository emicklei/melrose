package main

// script to generate the snippets.json for melrose-for-vscode

import (
	"encoding/json"
	"os"

	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
)

// see Makefile how to run this

type Snippet struct {
	Prefix      string   `json:"prefix"`
	Body        []string `json:"body"`
	Description string   `json:"description"`
}

func snippets() {
	snippets := map[string]Snippet{}
	ctx := melrose.PlayContext{
		VariableStorage: dsl.NewVariableStore(),
		LoopControl:     melrose.NoLooper,
	}
	for _, v := range dsl.EvalFunctions(ctx) {
		if len(v.Prefix) > 0 && len(v.Title) > 0 {
			snip := Snippet{
				Prefix:      v.Prefix,
				Body:        []string{v.Template},
				Description: v.Description,
			}
			snippets[v.Title] = snip
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(snippets)
}
