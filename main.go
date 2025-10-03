package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed instructions.md
var instructions string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	run := mcp.NewTool("ast-grep",
		mcp.WithDescription("Run ast-grep commands. ast-grep is a very fast CLI tool for code structural search, lint and rewriting.`"),
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

	s := server.NewMCPServer(
		"mcp-ast-grep",
		"1.0.2",
		server.WithInstructions(instructions),
	)
	s.AddTools(
		server.ServerTool{Tool: run, Handler: runHandler},
		server.ServerTool{Tool: help, Handler: helpHandler},
	)

	srv := server.NewStdioServer(s)
	fmt.Fprintln(os.Stderr, "mcp-ast-grep server started")
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
