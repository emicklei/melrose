package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/emicklei/melrose/api"
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPServer struct {
	service api.Service
}

func NewMCPServer(ctx core.Context) *MCPServer {
	// do not write to stdout as the MCP server is using that
	notify.Console.StandardOut = os.Stderr
	return &MCPServer{service: api.NewService(ctx)}
}

func (s *MCPServer) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	expression, ok := request.Params.Arguments["expression"].(string)
	if !ok {
		return nil, errors.New("expression must be a string")
	}
	toolResult := new(mcp.CallToolResult)
	result, err := s.service.CommandPlay("melrose-mcp", 0, expression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "play failed: %v\n", err)
		toolResult.IsError = true
		toolResult.Content = []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			}}
		return toolResult, err
	}
	duration := ""
	if endTime, ok := result.(time.Time); ok {
		duration = time.Until(endTime).String()
	}
	toolResult.Content = []mcp.Content{
		mcp.TextContent{
			Type: "text",
			Text: duration,
		},
	}
	return toolResult, nil
}
