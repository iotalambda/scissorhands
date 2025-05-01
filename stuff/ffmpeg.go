package stuff

import (
	"fmt"
	"os/exec"
)

func FfmpegExtractAudio(inputFilePath string, outputFilePath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-acodec", "copy", outputFilePath, "-y")
	// cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("extract audio with ffmpeg: %v", err)
	}
	return nil
}
