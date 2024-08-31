package models

import (
	"fmt"
	"runtime"

	"github.com/LuizGuilherme13/git-clone/common"
)

func CheckError(msg string) error {
	_, file, line, _ := runtime.Caller(1)

	return fmt.Errorf("%sError:%s %s (at %s:%d)", common.ColorRed, common.ColorReset, msg, file, line)
}

func PrintOk(msg string) {
	fmt.Printf("%s %s\n", common.Okay, msg)
}
