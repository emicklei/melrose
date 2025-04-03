package mcpserver

import (
	"context"
	"errors"
	"time"

	"github.com/emicklei/melrose/api"
	"github.com/emicklei/melrose/core"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPServer struct {
	service api.Service
}

func NewMCPServer(ctx core.Context) *MCPServer {
	return &MCPServer{service: api.NewService(ctx)}
}

func (s *MCPServer) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	expression, ok := request.Params.Arguments["expression"].(string)
	if !ok {
		return nil, errors.New("expression must be a string")
	}
	toolResult := new(mcp.CallToolResult)
	result, err := s.service.CommandPlay("melrose-mcp", 0, expression)
	if endTime, ok := result.(time.Time); ok {
		time.Sleep(time.Until(endTime))
	}
	if err != nil {
		toolResult.IsError = true
		toolResult.Content = []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			}}
		return toolResult, err
	}
	toolResult.Content = []mcp.Content{
		mcp.TextContent{
			Type: "text",
			Text: core.Storex(result),
		},
	}
	return toolResult, nil
}
