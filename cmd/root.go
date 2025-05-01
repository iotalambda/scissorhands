package cmd

import (
	"fmt"
	"log"
	"scissorhands/config"
	"strings"

	"github.com/spf13/cobra"
)

var input string
var maxSpeakers int
var message string
var output string
var service string

var rootCmd = &cobra.Command{
	Use:   "scissorhands",
	Short: "A CLI tool for smart video editing",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.InitConfig()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func flagsRequiredError(flags ...string) error {
	return fmt.Errorf("flags required: %v", strings.Join(flags, ", "))
}
