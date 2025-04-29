package main

import (
	"context"
	"fmt"
	"os"

	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/melrose/server/mcpserver"
	"github.com/emicklei/melrose/system"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var BuildTag = "dev"

func main() {
	notify.SetANSIColorsEnabled(false) // error messages cannot be colored

	ctx, err := system.Setup(BuildTag)
	if err != nil {
		notify.Errorf("setup failed: %v", err)
		os.Exit(1)
	}

	ioServer := server.NewMCPServer(
		"melrose",
		"v0.56.0",
	)
	playServer := mcpserver.NewMCPServer(ctx)

	// Add play tool
	tool := mcp.NewTool("melrose-play",
		mcp.WithDescription("play melrōse expression"),
		mcp.WithString("expression",
			mcp.Required(),
			mcp.Description("melrōse expression to play"),
		),
	)
	ioServer.AddTool(tool, playServer.HandlePlay)

	// Add chord prompt
	chordHander := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		note := request.Params.Arguments["ground"]
		if note == "" {
			note = "C"
		}
		fraction := request.Params.Arguments["fraction"]
		if fraction == "" {
			fraction = "4"
		}
		octave := request.Params.Arguments["octave"]
		if octave == "" {
			octave = "4"
		}
		return mcp.NewGetPromptResult(
			"playing a chord",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent(fmt.Sprintf("chord('%s%s%s')", fraction, note, octave)),
				),
			},
		), nil
	}
	ioServer.AddPrompt(mcp.NewPrompt("play-chord",
		mcp.WithPromptDescription("play the notes of a chord")), chordHander)

	// Start the stdio server
	if err := server.ServeStdio(ioServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
