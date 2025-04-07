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
		mcp.WithDescription("play melrōse expression"),
		mcp.WithString("expression",
			mcp.Required(),
			mcp.Description("melrōse expression to play"),
		),
	)

	playServer := mcpserver.NewMCPServer(ctx)

	// Add tool handler
	ioServer.AddTool(tool, playServer.Handle)

	/**
		ioServer.AddPrompt(mcp.NewPrompt("play-a-chord"),
			mcp.WithPromptDescription("play the notes of a chord"),
			mcp.WithArgument("ground",
				mcp.ArgumentDescription("the ground note of the chord"),
			), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				name := request.Params.Arguments["ground"]
				if name == "" {
					name = "ground"
				}

				return mcp.NewGetPromptResult(
					"A friendly greeting",
					[]mcp.PromptMessage{
						mcp.NewPromptMessage(
							mcp.RoleAssistant,
							mcp.NewTextContent(fmt.Sprintf("Hello, %s! How can I help you today?", name)),
						),
					},
				), nil
			})
	**/

	// Start the stdio server
	if err := server.ServeStdio(ioServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
