package main

import (
	"fclip"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/phuslu/log"
	"os"
)
var cli struct {
	Files []string `arg type:"existingfile"`
}

func main() {
	initLog()
	// no args
	if len(os.Args) == 1 {
		// be silent
		log.DefaultLogger.Level = log.PanicLevel
		printClipboardFiles()
		return
	}
	//fclip.PathsToClipboard(`c:\1\TessaTalks.pdf`)
/*	paths, err := fclip.GetPathsFromClipboard()
	if err != nil {
		//panic(err)
	}
	fmt.Printf("Clipboard data: %v", paths)
*/

	ctx := kong.Parse(&cli,
		kong.Name("fclip"),
		kong.Description("Simple utility for adding or getting files to and from Windows clipboard."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	log.Info().Str("command", ctx.Command())
	switch ctx.Command() {}
}

func initLog() {
	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			Level: log.TraceLevel,
			TimeFormat: "15:04:05",
			Caller:     1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}
}

func printClipboardFiles() {
	log.Info().Msg("print clipboard files")
	paths, err := fclip.GetPathsFromClipboard()
	if err != nil {
		panic(err)
	}
	for _, p := range paths {
		fmt.Println(p)
	}
}