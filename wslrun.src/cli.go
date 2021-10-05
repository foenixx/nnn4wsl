package main

import "github.com/alecthomas/kong"


type CLI struct {
	Gui bool `help:"Run in GUI mode."`
	NoGui bool `help:"Run in console mode."`
	Verbose bool `help:"Print logs."`
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
