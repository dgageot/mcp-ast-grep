# MCP Server for ast-grep

[ast-grep](https://github.com/ast-grep/ast-grep) is a CLI tool for code structural search, lint and rewriting. Written in Rust.

## Usage with docker mcp

```sh
docker mcp tools call ast-grep args=run args="--pattern" args='const $NAME = $VAL' args="--json" args="."
```