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
		mcp.WithDescription("Search codebases using ast-grep patterns for straightforward structural matches. (returns json)"),
		mcp.WithTitleAnnotation("Search codebases using ast-grep patterns"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("pattern",
			mcp.Description("The ast-grep pattern to search for. Note, the pattern must have valid AST structure."),
			mcp.Required(),
		),
		mcp.WithString("lang",
			mcp.Description("The language of the code. If not specified, will be auto-detected based on file extensions."),
		),
		mcp.WithString("dir",
			mcp.Description("The directory to search in. If not specified, will use the current working directory."),
		),
	)

	s := server.NewMCPServer(
		"mcp-ast-grep",
		"1.0.2",
		server.WithInstructions(instructions),
	)
	s.AddTools(
		server.ServerTool{Tool: run, Handler: runHandler},
	)

	srv := server.NewStdioServer(s)
	fmt.Fprintln(os.Stderr, "mcp-ast-grep server started")
	if err := srv.Listen(ctx, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func runHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pattern, err := request.RequireString("pattern")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"run", "--pattern", pattern, "--json"}
	if lang := request.GetString("lang", ""); lang != "" {
		args = append(args, "--lang", lang)
	}
	if dir := request.GetString("dir", "."); dir != "" {
		args = append(args, dir)
	}

	return runAstGrep(ctx, args...), nil
}

func runAstGrep(ctx context.Context, args ...string) *mcp.CallToolResult {
	out, err := exec.CommandContext(ctx, "/usr/bin/ast-grep", args...).CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Command failed: %s\nOutput: %s", err.Error(), string(out)))
	}

	return mcp.NewToolResultText(string(out))
}
