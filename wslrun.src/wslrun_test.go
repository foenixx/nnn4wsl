package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_getCmdPath(t *testing.T) {
	p, a := getCmdPathArgs(`C:\Users\oleg\AppData\Local\Programs\Microsoft VS Code\Code.exe`)
	assert.Equal(t, `C:\Users\oleg\AppData\Local\Programs\Microsoft`, p)
	assert.Equal(t, `VS Code\Code.exe`, a)
	p, a = getCmdPathArgs(`"C:\Users\oleg\AppData\Local\Programs\Microsoft VS Code\Code.exe"`)
	assert.Equal(t, `C:\Users\oleg\AppData\Local\Programs\Microsoft VS Code\Code.exe`, p)
	assert.Equal(t, "", a)
	p, a = getCmdPathArgs(`do.exe param1 param2="value with space"`)
	assert.Equal(t, `do.exe`, p)
	assert.Equal(t, `param1 param2="value with space"`, a)
	p, a = getCmdPathArgs(`'./test 1.bat' param1 param2`)
	assert.Equal(t, `./test 1.bat`, p)
	assert.Equal(t, `param1 param2`, a)
}

func Test_wslPath2Win(t *testing.T) {
	assert.Equal(t, "c:/some/path", wslPath2Win("/mnt/c/some/path"))
	assert.Equal(t, `c:\some\path`, wslPath2Win(`c:\some\path`))
	os.Setenv("WSLR_UNC_PREFIX", `\\wsl$\Ubuntu`)
	assert.Equal(t, `\\wsl$\Ubuntu\some\path`, wslPath2Win(`/some/path`))
	assert.Equal(t, `some_path`, wslPath2Win(`some_path`))
	assert.Equal(t, `./some_path`, wslPath2Win(`./some_path`))
}