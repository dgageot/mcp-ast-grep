package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const astGrepPath = "/usr/bin/ast-grep"

func main() {
	s := server.NewMCPServer("mcp-ast-grep", "1.0.0")

	run := mcp.NewTool("ast-grep",
		mcp.WithDescription("Run ast-grep commands. ast-grep is a very fast CLI tool for code structural search, lint and rewriting."),
		mcp.WithTitleAnnotation("Run ast-grep"),
		mcp.WithArray("args",
			mcp.Required(),
			mcp.Description("Arguments for the ast-grep command"),
			mcp.WithStringItems(),
		),
	)
	s.AddTool(run, runHandler)

	help := mcp.NewTool("ast-grep-help",
		mcp.WithDescription("Get help for ast-grep commands."),
		mcp.WithTitleAnnotation("Get help for ast-grep"),
	)
	s.AddTool(help, helpHandler)

	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

func runHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := request.RequireStringSlice("args")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	out, err := exec.CommandContext(ctx, astGrepPath, args...).CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Command failed: %s\nOutput: %s", err.Error(), string(out))), nil
	}

	return mcp.NewToolResultText(string(out)), nil
}

func helpHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	out, err := exec.CommandContext(ctx, astGrepPath, "--help").CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Command failed: %s\nOutput: %s", err.Error(), string(out))), nil
	}

	return mcp.NewToolResultText(string(out)), nil
}
