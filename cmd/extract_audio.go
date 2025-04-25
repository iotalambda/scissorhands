package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var extractAudioCmd = &cobra.Command{
	Use:   "extract-audio",
	Short: "Extract audio from an input file",
	Run: func(_ *cobra.Command, _ []string) {
		if err := extractAudio(); err != nil {
			panic(err)
		}
	},
}

func extractAudio() error {
	cmd := exec.Command("ffmpeg", "-i", input, "-vn", "-acodec", "copy", output, "-y")
	// cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("extract audio with ffmpeg: %v", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(extractAudioCmd)

	extractAudioCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	extractAudioCmd.MarkFlagRequired("input")

	extractAudioCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path.")
	extractAudioCmd.MarkFlagRequired("output")
}
