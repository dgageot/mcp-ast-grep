package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	run := mcp.NewTool("ast-grep",
		mcp.WithDescription("Run ast-grep commands. ast-grep is a very fast CLI tool for code structural search, lint and rewriting. An example of usage is `ast-grep run -l go --pattern 'const $NAME = $VAL' --json .`"),
		mcp.WithTitleAnnotation("Run ast-grep"),
		mcp.WithArray("args",
			mcp.Description("Arguments for the ast-grep command"),
			mcp.WithStringItems(),
			mcp.Required(),
		),
	)
	help := mcp.NewTool("ast-grep-help",
		mcp.WithDescription("Get help for ast-grep commands."),
		mcp.WithTitleAnnotation("Get help for ast-grep"),
		mcp.WithReadOnlyHintAnnotation(true),
	)

	s := server.NewMCPServer("mcp-ast-grep", "1.0.1")
	s.AddTools(
		server.ServerTool{Tool: run, Handler: runHandler},
		server.ServerTool{Tool: help, Handler: helpHandler},
	)

	srv := server.NewStdioServer(s)
	if err := srv.Listen(ctx, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func runHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := request.RequireStringSlice("args")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return runAstGrep(ctx, args...), nil
}

func helpHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return runAstGrep(ctx, "--help"), nil
}

func runAstGrep(ctx context.Context, args ...string) *mcp.CallToolResult {
	out, err := exec.CommandContext(ctx, "/usr/bin/ast-grep", args...).CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Command failed: %s\nOutput: %s", err.Error(), string(out)))
	}

	return mcp.NewToolResultText(string(out))
}
