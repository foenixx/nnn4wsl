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

func getCmdPathArgs(cmd string) (string, string) {
	// "comm.exe " param1 param2
	// " c:\program files\command\comm.exe  "
	// comm.exe param1 param2
	var reCmdName = regexp.MustCompile(`^\s*(".+?"|'.+?'|\S+)\s*(.*)$`)

	m := reCmdName.FindStringSubmatch(cmd)
	if m == nil {
		return "", ""
	}
	if m[1][0] == '"' || m[1][0] == '\'' {
		//strip quotes
		m[1] = m[1][1 : len(m[1])-1]
	}

	return m[1], m[2]

}

func wslPath2Win(path string) string {
	if strings.HasPrefix(path, "/mnt/") {
		// /mnt/c/some/path --> c:/some/path
		return fmt.Sprintf("%s:%s", string(path[5]),path[6:])
	}
	// path is absolute?
	if strings.HasPrefix(path, `/`) {
		// append unc prefix like "\\wsl$\Ubuntu-20.04"
		return os.Getenv("WSLR_UNC_PREFIX") + strings.ReplaceAll(path, `/`,`\`)
	}
	return path
}

func pwshStartProcess(gui bool, process string, args string) error {
	var a = []string{"-Command", "Start-Process", "\"" + process + "\""}
	if !gui {
		a = append(a, "-NoNewWindow", "-Wait")
	}
	if args != "" {
		// Input:  wslrun go build -ldflags="-s -w"
		// Output: pwsh.exe -Command Start-Process go -ArgumentList "build -ldflags=`"-s -w`""
		l := `"` + strings.ReplaceAll(args, "\"", "`\"") + `"`
		a = append(a, "-ArgumentList", l)
	}
	log.Info().Strs("args", a).Msg("powershell process arguments")

	return pwsh(a...)
}

func pwsh(args ...string) error {
	pwsh, err := exec.LookPath("pwsh.exe")
	if err != nil {
		log.Info().Msg("cannot find pwsh.exe in the PATH, trying to use environment variable WSLR_PWSH")
		pwsh = os.Getenv("WSLR_PWSH")
		if pwsh == "" {
			log.Error().Msg("cannot find pwsh.exe, giving up...")
			os.Exit(1)
		}
	}
	cmd := exec.Command(pwsh, args...)
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
	} else {
		if len(cli.Argv) == 0 {
			log.Error().Msg("no execution parameters passed")
			os.Exit(1)
		}
		log.Info().Strs("argv", cli.Argv).Msg("exec arguments")
		buf = strings.Join(cli.Argv, " ")
	}
	path, args := getCmdPathArgs(buf)
	// trying to translate wsl path into win path
	path = wslPath2Win(path)
	ext := filepath.Ext(path)
	log.Info().Str("path", path).Str("args", args).Str("file ext", ext).Msg("command parameters")

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
	err = pwshStartProcess(gui, path, args)

	if err != nil {
		log.Error().Err(err).Msg("powershell launch error")
		os.Exit(1)
	}
}
