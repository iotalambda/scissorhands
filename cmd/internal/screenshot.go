package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var ss string
var targetFormat string

var ScreenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Take a screenshot from a video at a certain point in time",
	RunE: func(cmd *cobra.Command, args []string) error {
		return screenshot()
	},
}

func screenshot() error {
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
	ScreenshotCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	ScreenshotCmd.MarkFlagRequired("input")

	ScreenshotCmd.Flags().StringVarP(&targetFormat, "target-format", "t", "", "Target format of the output. Allowed values: base64, file.")

	ScreenshotCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path.")

	ScreenshotCmd.Flags().StringVarP(&ss, "seek-start", "s", "", "Seek start, for example \"00:59:59.999\".")
	ScreenshotCmd.MarkFlagRequired("ss")
}
