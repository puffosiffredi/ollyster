package conf

import (
	"bufio"
	"log"
	"ollyster/tools"
	"os"
	"path/filepath"
	"strings"
)

var OProfile map[string]string

func init() {

	OProfile = make(map[string]string)

	ProfileFile := filepath.Join(tools.Hpwd(), "etc", "profile.conf")

	readProfile(ProfileFile)

}

func StartProfile() {

	log.Printf("[ETC] Reading profile...\r\n")

}

func serializeProfile(line string) {

	// create a splitter because "split" adds an empty line after the last \n
	splitter := func(c rune) bool {
		return (c == ' ' || c == '=') // trims space and understands equal
	}

	split := strings.FieldsFunc(line, splitter)

	if len(split) != 0 {

		OProfile[split[0]] = split[1]
		log.Printf("[ETC] Profile: %q -> %q\r\n", split[0], split[1])

	}

}

func readProfile(FileName string) {

	file, err := os.Open(FileName)
	if err != nil {
		log.Printf("[ETC] can't open file %s", FileName)

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 4 {
			serializeProfile(line)
		}
	}

	file.Close()

}
