package plugin

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

func printOut(list appItems) {
	if !list.hasChanged {
		return
	}
	// jsonBytes, err := json.MarshalIndent(list.data, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("=====", string(jsonBytes))

	// show as yaml
	yamlData, err := yaml.Marshal(list.data)
	if err != nil {
		panic(err)
	}
	// fmt.Println(">>>>>")
	fmt.Println(string(yamlData))
	fmt.Println("---")
}

func splitasPath(value string) ([]string, string) {
	fulllist := strings.Split(value, ".")

	// fmt.Println(">>", fulllist)

	if len(fulllist) == 0 {
		return []string{}, value
	}
	if len(fulllist) == 1 {
		return []string{}, value
	}

	// fmt.Println("!>", fulllist[:len(fulllist)-1])

	return fulllist[:len(fulllist)-1], fulllist[len(fulllist)-1]
}
