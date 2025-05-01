package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var azureSpeechDiarizationInput string
var openAIWhisperSegmentationInput string

var createDialogueCmd = &cobra.Command{
	Use:   "create-dialogue",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return createDialogue()
	},
}

func createDialogue() error {
	cmdArgs := []string{"-ss", ss, "-i", input, "-vframes", "1"}
	switch targetFormat {
	case "base64":
		cmdArgs = append(cmdArgs, "-f", "image2", "-")
	case "file":
		if output == "" {
			return flagsRequiredError("--output")
		}
		cmdArgs = append(cmdArgs, output)
	default:
		return fmt.Errorf("unrecognized target format: %v", targetFormat)
	}
	cmd := exec.Command("ffmpeg", cmdArgs...)
	// cmd.Stderr = os.Stderr

	var out bytes.Buffer
	if targetFormat == "base64" {
		cmd.Stdout = &out
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("screenshot with ffmpeg: %v", err)
	}

	if targetFormat == "base64" {
		url := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(out.Bytes())
		fmt.Println(url)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(createDialogueCmd)

	createDialogueCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	createDialogueCmd.MarkFlagRequired("input")

	createDialogueCmd.Flags().StringVarP(&targetFormat, "target-format", "t", "", "Target format of the output. Allowed values: base64, file.")

	createDialogueCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path.")

	createDialogueCmd.Flags().StringVarP(&ss, "seek-start", "s", "", "Seek start, for example \"00:59:59.999\".")
	createDialogueCmd.MarkFlagRequired("ss")
}
