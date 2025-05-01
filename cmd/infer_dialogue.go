package cmd

import (
	"github.com/spf13/cobra"
)

var inferDialogueCmd = &cobra.Command{
	Use:   "infer-dialogue",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return inferDialogue()
	},
}

func inferDialogue() error {
	return nil
}

func init() {

}
