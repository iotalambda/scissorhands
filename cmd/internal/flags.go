package internal

import (
	"fmt"
	"strings"
)

var input string
var maxSpeakers int
var message string
var output string
var service string

func flagsRequiredError(flags ...string) error {
	return fmt.Errorf("flags required: %v", strings.Join(flags, ", "))
}
