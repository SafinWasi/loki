package tools

import (
	"log"

	"github.com/safinwasi/loki/cmd"
)

func DebugPrint(input string) {
	if cmd.Debug {
		log.Println(input)
	}
}
