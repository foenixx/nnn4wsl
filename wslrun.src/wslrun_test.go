package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getCmdPath(t *testing.T) {
	p,a := getCmdPathArgs(`C:\Users\oleg\AppData\Local\Programs\Microsoft VS Code\Code.exe`)
	assert.Equal(t, `C:\Users\oleg\AppData\Local\Programs\Microsoft` , p)
	assert.Equal(t, `VS Code\Code.exe`, a)
	p, a = getCmdPathArgs(`"C:\Users\oleg\AppData\Local\Programs\Microsoft VS Code\Code.exe"`)
	assert.Equal(t, `C:\Users\oleg\AppData\Local\Programs\Microsoft VS Code\Code.exe`, p)
	assert.Equal(t, "", a)
	p, a = getCmdPathArgs(`do.exe param1 param2="value with space"`)
	assert.Equal(t, `do.exe`, p)
	assert.Equal(t, `param1 param2="value with space"`, a)
}
