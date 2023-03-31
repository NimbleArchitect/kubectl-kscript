package plugin

import (
	"fmt"
	"strings"
	"time"
)

type globalObject struct {
	Namespace string
	Cluster   string
	Selector  string
}

type jsConsole struct {
	start time.Time
}

func (d *jsConsole) Log(s ...string) {
	fmt.Println(strings.Join(s, " "))
}
