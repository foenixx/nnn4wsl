package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)


func getCmdPathArgs(cmd string) (string, string) {
	// "comm.exe " param1 param2
	// " c:\program files\command\comm.exe  "
	// comm.exe param1 param2
	var reCmdName = regexp.MustCompile(`^(?:"\s*([\S ]*?)\s*"|(\S+))\s*(.*)$`)

	m := reCmdName.FindStringSubmatch(cmd)
	if m == nil {
		return "", ""
	}
	var args string
	if len(m) == 4 {
		// there are command arguments in the cmd
		args = m[3]
	}
	if m[1] != "" {
		// double-quoted command path
		return m[1], args
	}
	// unquoted command
	return m[2], args
}

func pwshStartProcess(gui bool, process string, args string) error {
	var a = []string {"-Command", "Start-Process", "\"" + process + "\""}
	if !gui {
		a = append(a, "-NoNewWindow", "-Wait")
	}
	if args != "" {
		// Input:  wslrun go build -ldflags="-s -w"
		// Output: pwsh.exe -Command Start-Process go -ArgumentList "build -ldflags=`"-s -w`""
		l := `"` + strings.ReplaceAll(args, "\"", "`\"") + `"`
		a = append(a,"-ArgumentList", l)
	}
	log.Printf("pwsh rguments: %s", a)
	return pwsh(a...)
}

func pwsh(args... string) error {
	pwsh, err := exec.LookPath( "pwsh.exe" )
	if err != nil {
		log.Printf("cannot find pwsh.exe in the PATH, trying to use environment variable WSLR_PWSH")
		pwsh = os.Getenv("WSLR_PWSH")
		if pwsh == "" {
			log.Printf("cannot find pwsh.exe, giving up...")
			os.Exit(1)
		}
	}
	cmd := exec.Command(pwsh, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("PowerShell error: %s", err)
		return err
	}
	return nil
}

func main() {
	cli := cliParse()
	if !cli.Verbose {
		log.SetOutput(ioutil.Discard)
	}
	buf := os.Getenv("WSLR_BUFFER")
	if buf == "" {
		log.Print("error: empty WSLR_BUFFER")
		os.Exit(1)
	}
	log.Printf("buffer: %s", buf)
	path,args := getCmdPathArgs(buf)
	ext := filepath.Ext(path)
	log.Printf("command path: %s", path)
	log.Printf("command args: %s", args)
	log.Printf("command ext: %s", ext)

	var err error
	var gui bool
	switch ext {
	case ".exe", "bat":
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
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}
