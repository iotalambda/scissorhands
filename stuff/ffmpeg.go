package stuff

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func FfmpegExtractAudio(inputFilePath string, outputFilePath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-acodec", "copy", outputFilePath, "-y")
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg extract audio: %v", err)
	}
	return nil
}

func FfprobeDurationMs(inputFilePath string) (int, error) {
	cmd := exec.Command("ffprobe", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", inputFilePath)
	// cmd.Stderr = os.Stderr
	outBytes, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe duration: %v", err)
	}

	outStr := strings.Split(string(outBytes), "\n")[0]
	outFloatS, err := strconv.ParseFloat(outStr, 32)
	if err != nil {
		return 0, fmt.Errorf("parse float: %v", err)
	}

	outIntMs := int(outFloatS * 1000)
	return outIntMs, nil
}

func MsToSeek(ms int) string {
	msTgt := ms % 1000
	sTgt := (ms / 1000) % 60
	mTgt := (ms / 1000 / 60) % 60
	hTgt := (ms / 1000 / 60 / 60)
	ssTgt := fmt.Sprintf("%02d:%02d:%02d.%03d", hTgt, mTgt, sTgt, msTgt)
	return ssTgt
}

func FfmpegScreenshot(ss string, inputFilePath string) (string, error) {
	cmd := exec.Command("ffmpeg", "-ss", ss, "-i", inputFilePath, "-vframes", "1", "-f", "image2", "-")
	// cmd.Stderr = os.Stderr
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg: %v", err)
	}
	url := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(out.Bytes())
	return url, nil
}
