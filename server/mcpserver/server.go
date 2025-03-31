package mcpserver

import (
	"context"
	"errors"
	"strings"

	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/dsl"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPServer struct {
	context   core.Context
	evaluator *dsl.Evaluator
}

func NewMCPServer(ctx core.Context) *MCPServer {
	return &MCPServer{context: ctx, evaluator: dsl.NewEvaluator(ctx)}
}

func (s *MCPServer) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	expression, ok := request.Params.Arguments["expression"].(string)
	if !ok {
		return nil, errors.New("expression must be a string")
	}
	var playExpression string
	if strings.HasPrefix(expression, "play(") {
		playExpression = expression
	} else {
		playExpression = "play(" + expression + ")"
	}
	result, err := s.evaluator.EvaluateExpression(playExpression)
	if err != nil {
		return mcp.NewToolResultText(err.Error()), err
	}
	inspect := core.NewInspect(s.context, "result", result)
	return mcp.NewToolResultText(inspect.Markdown()), nil
}
