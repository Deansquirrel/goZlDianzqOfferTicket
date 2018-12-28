package common

import (
	"fmt"
	"github.com/Deansquirrel/go-tool"
)

var IsDebug bool

func MyLog(s string) {
	if IsDebug {
		fmt.Println(s)
	} else {
		err := go_tool.Log(s)
		if err != nil {
			fmt.Println(err, " - ", s)
		}
	}
}
