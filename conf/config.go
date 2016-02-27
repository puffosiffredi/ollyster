package conf

import (
	"bufio"
	"log"
	"ollyster/tools"
	"os"
	"path/filepath"
	"strings"
)

var OConfig map[string]string

func init() {

	OConfig = make(map[string]string)

	GConfigFile := filepath.Join(tools.Hpwd(), "etc", "ollyster.conf")

	readConfig(GConfigFile)

}

func StartConfig() {

	log.Printf("[ETC] Reading config...\r\n")

}

func serializeConf(line string) {

	// create a splitter because "split" adds an empty line after the last \n
	splitter := func(c rune) bool {
		return (c == ' ' || c == '=') // trims space and understands equal
	}

	split := strings.FieldsFunc(line, splitter)

	if len(split) != 0 {

		OConfig[split[0]] = split[1]
		log.Printf("[ETC] Config: %q -> %q\r\n", split[0], split[1])

	}

}

func readConfig(FileName string) {

	file, err := os.Open(FileName)
	if err != nil {
		log.Printf("[ETC] can't open file %s", FileName)

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 4 {
			serializeConf(line)
		}
	}

	file.Close()

}

func GetConfItem(item_name string) string {

	var tmp_value string

	if val, ok := OConfig[item_name]; ok {

		tmp_value = val
	} else {

		log.Printf("[ETC][ERROR] no config item named %s", item_name)

		tmp_value = "ERROR"
	}

	return tmp_value

}
