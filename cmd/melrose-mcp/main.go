package main

import (
	"context"
	"fmt"
	"log"

	"github.com/emicklei/melrose/server/mcpserver"
	"github.com/emicklei/melrose/system"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var BuildTag = "dev"

func main() {
	ctx, err := system.Setup(BuildTag)
	if err != nil {
		log.Fatalln(err)
	}

	ioServer := server.NewMCPServer(
		"melrōse",
		"v0.56.0",
	)

	// Add tool
	tool := mcp.NewTool("play-melrose",
		mcp.WithDescription("play melrōse expression"),
		mcp.WithString("expression",
			mcp.Required(),
			mcp.Description("melrōse expression to play"),
		),
	)

	playServer := mcpserver.NewMCPServer(ctx)

	// Add tool handler
	ioServer.AddTool(tool, playServer.Handle)

	chordHander := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		note := request.Params.Arguments["ground"]
		if note == "" {
			note = "ground"
		}
		return mcp.NewGetPromptResult(
			"playing a chord",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent(fmt.Sprintf("chord('%s')", note)),
				),
			},
		), nil
	}
	chordPrompt := mcp.NewPrompt("play-a-chord",
		mcp.WithPromptDescription("play the notes of a chord"))

	ioServer.AddPrompt(chordPrompt, chordHander)

	// Start the stdio server
	if err := server.ServeStdio(ioServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
