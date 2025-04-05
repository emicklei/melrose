package main

import (
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
		mcp.WithDescription("play note sequences"),
		mcp.WithString("expression",
			mcp.Required(),
			mcp.Description("melrōse expression to play"),
		),
	)

	playServer := mcpserver.NewMCPServer(ctx)

	// Add tool handler
	ioServer.AddTool(tool, playServer.Handle)

	// ioServer.AddPrompt(mcp.NewPrompt("chord"),
	// 	mcp.WithPromptDescription("chord with 3 notes"),
	// 	mcp.WithArgument("base", mcp.With

	// )

	// Start the stdio server
	if err := server.ServeStdio(ioServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
