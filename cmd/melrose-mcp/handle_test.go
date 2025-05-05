package main

import (
	"context"
	"log"
	"testing"

	"github.com/emicklei/melrose/server/mcpserver"
	"github.com/emicklei/melrose/system"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleCDE(t *testing.T) {
	t.Skip()
	ctx, err := system.Setup("test")
	if err != nil {
		log.Fatalln(err)
	}
	defer ctx.Device().Close()
	playServer := mcpserver.NewMCPServer(ctx)

	req := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "play-melrose",
		},
	}
	req.Params.Name = "play-melrose"
	req.Params.Arguments = map[string]interface{}{
		"expression": `a=note('c')
b=a+a`,
	}
	result, err := playServer.HandlePlay(context.Background(), req)
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("Handle returned error: %v", result)
	}
	t.Log("Handle result:", result.Result)
}
