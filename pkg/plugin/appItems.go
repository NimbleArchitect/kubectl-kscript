package plugin

import (
	"strconv"
	"strings"
)

type appItems struct {
	hasChanged bool
	data       map[string]interface{}
}

func (l *appItems) add(name string, data interface{}, path []string) {
	if l == nil {
		return
	}

	if data == nil {
		return
	}

	if l.data == nil {
		l.data = make(map[string]interface{})
	}

	list := l.data
	for _, part := range path {
		if strings.Contains(part, "[") {
			indexStr := part[strings.Index(part, "[")+1 : strings.Index(part, "]")]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				panic(err)
			}
			part = part[:strings.Index(part, "[")]
			if list[part] == nil {
				list[part] = make([]interface{}, index+1)

			} else {
				newArray := make([]interface{}, index+1)
				copy(newArray, list[part].([]interface{}))
				list[part] = newArray
			}

			array := list[part].([]interface{})
			if array[index] == nil {
				array[index] = make(map[string]interface{})
			}
			list = array[index].(map[string]interface{})
		} else {
			if list[part] == nil {
				list[part] = make(map[string]interface{})
			}
			list = list[part].(map[string]interface{})
		}
	}

	list[name] = data
	l.hasChanged = true
}
