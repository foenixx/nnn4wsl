package main

import "github.com/alecthomas/kong"

type CLI struct {
	Env     bool     `help:"Use WSLR_BUFFER variable as input string and ignore command-line exec arguments." short:"e"`
	Gui     bool     `help:"Run in GUI mode."`
	NoGui   bool     `help:"Run in console mode."`
	Hidden  bool     `help:"Run with hidden console window (hack for VS Code)."`
	Verbose bool     `help:"Print verbose logs. Also you can set WSLR_VERBOSE=1 environment variable to switch verbose output on."`
	Argv    []string `arg:"" optional:"" passthrough:""`
}

func cliParse() *CLI {
	var cli CLI
	kong.Parse(&cli,
		kong.Name("wslrun"),
		kong.Description("An util for running windows commands and executables from WSL2"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	return &cli
}
