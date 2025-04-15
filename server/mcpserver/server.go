package mcpserver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
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
	return &MCPServer{service: api.NewService(ctx)}
}

func (s *MCPServer) HandlePlay(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	expression, ok := request.Params.Arguments["expression"].(string)
	if !ok {
		return nil, errors.New("expression must be a string")
	}
	toolResult := new(mcp.CallToolResult)

	// do not write to stdout as the MCP server is using that
	captured := new(bytes.Buffer)
	notify.Console.StandardOut = captured

	response, err := s.service.CommandPlay("melrose-mcp", 0, expression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "play failed: %v\n", err)
		toolResult.IsError = true
		toolResult.Content = []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: expression,
			},
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			}}
		return toolResult, err
	}
	content := []mcp.Content{
		mcp.TextContent{
			Type: "text",
			Text: time.Until(response.EndTime).String(),
		}, mcp.TextContent{
			Type: "text",
			Text: strings.TrimSpace(captured.String()),
		}}
	if p, ok := response.ExpressionResult.(core.Sequenceable); ok {
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: p.S().Storex(),
		})
	} else {
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("%v", response.ExpressionResult),
		})
	}
	toolResult.Content = content
	return toolResult, nil
}
