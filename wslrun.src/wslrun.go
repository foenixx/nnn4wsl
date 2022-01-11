package main

import (
	"fmt"
	"github.com/phuslu/log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func initLog(verbose bool) {
	//if log.IsTerminal(os.Stderr.Fd()) {
	w := &log.ConsoleWriter{
		ColorOutput:    true,
		QuoteString:    false,
		EndWithMessage: false,
	}
	//}
	var lvl = log.ErrorLevel
	env := os.Getenv("WSLR_VERBOSE")
	if verbose || env != "" {
		lvl = log.DebugLevel
	}
	log.DefaultLogger = log.Logger{
		Level:      lvl,
		TimeFormat: "2006.01.02 15:04:05",
		Caller:     0,
		Writer:     w,
	}
}

// splitArgs splits argument list string into separate arguments.
// For example, this string:
//   "text w spaces" 'text w spaces'  text"text" "text"text text'text' 'text'text text w spaces
// splits into args:
//	"text w spaces",'text w spaces',text"text","text"text,text'text','text'text,text,w,spaces
func splitArgs(cmd string) []string {
	var reArgs = regexp.MustCompile(`(?:"(?:\\"|.)+?"\S*)|(?:'.+?'\S*)|(?:\S+)`)
	var reQuotes = regexp.MustCompile(`^(".*")|('.*')$`)
	var args []string
	matches := reArgs.FindAllString(cmd, -1)
	log.Info().Strs("matches", matches).Msg("parsing input")
	args = make([]string, len(matches))
	for i, m := range matches {
		if reQuotes.MatchString(m) {
			//split quotes
			args[i] = m[1 : len(m)-1]
		} else {
			args[i] = m
		}
	}
	log.Info().Strs("args", args).Msg("parsed")
	return args
}

func wslPath2Win(path string) string {
	if strings.HasPrefix(path, "/mnt/") {
		// /mnt/c/some/path --> c:/some/path
		return fmt.Sprintf("%s:%s", string(path[5]), path[6:])
	}
	// path is absolute?
	if strings.HasPrefix(path, `/`) {
		// append unc prefix like "\\wsl$\Ubuntu-20.04"
		return os.Getenv("WSLR_UNC_PREFIX") + strings.ReplaceAll(path, `/`, `\`)
	}
	return path
}

func pwshStartProcess(gui bool, hidden bool, process string, args []string) error {
	var a = []string{"-Command", "Start-Process", "'" + process + "'"}
	if !gui {
		a = append(a, "-NoNewWindow", "-Wait")
	}
	if hidden {
		// Visual Studio Code fix
		// https://mypowershellnotes.wordpress.com/2020/04/26/i-cant-close-my-shell-until-i-close-visual-studio-code/
		// https://stackoverflow.com/questions/57335812/how-to-open-visual-studio-code-through-powershell-and-close-powershell-right-aft
		a = append(a, "-WindowStyle Hidden")
	}

	if len(args) > 0 {
		// Input:  wslrun go build -ldflags="-s -w"
		// Output: pwsh.exe -Command Start-Process 'go' -ArgumentList 'build','-ldflags="-s -w"'
		a = append(a, "-ArgumentList '"+args[0]+"'")
		for _, arg := range args[1:] {
			a = append(a, ",'"+arg+"'")
		}
	}
	//a = append(a, "| Out-Null")
	log.Info().Strs("args", a).Msg("powershell process arguments")

	return pwsh(a...)
}

func pwsh(args ...string) error {
	exe, err := exec.LookPath("pwsh.exe")
	if err != nil {
		log.Info().Msg("cannot find pwsh.exe in the PATH, trying to use environment variable WSLR_PWSH")
		exe = os.Getenv("WSLR_PWSH")
		if exe == "" {
			log.Error().Msg("cannot find pwsh.exe, giving up...")
			os.Exit(1)
		}
	}
	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("PowerShell error")
		return err
	}
	return nil
}

func main() {
	var buf string
	var args []string
	cli := cliParse()
	initLog(cli.Verbose)
	if cli.Env {
		buf = os.Getenv("WSLR_BUFFER")
		if buf == "" {
			//log.Print("error: empty WSLR_BUFFER")
			log.Error().Msg("empty WSLR_BUFFER")
			os.Exit(1)
		}
		log.Info().Str("buf", buf).Msg("WSLR_BUFFER value")
		args = splitArgs(buf)
	} else {
		if len(cli.Argv) == 0 {
			log.Error().Msg("no execution parameters passed")
			os.Exit(1)
		}
		log.Info().Strs("argv", cli.Argv).Msg("exec arguments")
		args = cli.Argv
	}

	// trying to translate wsl path into win path
	path := wslPath2Win(args[0])
	ext := filepath.Ext(path)
	log.Info().Str("path", path).Strs("args", args[1:]).Str("file ext", ext).Msg("command parameters")

	var err error
	var gui bool
	switch ext {
	case ".exe", ".bat":
		gui = false
	case ".doc", ".docx", ".xls", ".xlsx", ".pdf", ".jpg", ".jpeg", ".png", ".sql", ".txt", ".lnk":
		gui = true
	}
	if cli.Gui {
		//force gui mode
		gui = true
	}
	if cli.NoGui {
		//force console mode
		gui = false
	}
	err = pwshStartProcess(gui, cli.Hidden, path, args[1:])

	if err != nil {
		log.Error().Err(err).Msg("powershell launch error")
		os.Exit(1)
	}
}
