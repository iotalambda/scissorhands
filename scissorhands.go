package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	if err := LoadConfig(); err != nil {
		return err
	}

	op, err := getOp()
	if err != nil {
		return err
	}

	switch op {
	case "extract-audio":
		return extractAudio()
	case "segment":
		return segment()
	default:
		return fmt.Errorf("operation `%v` not recognized", op)
	}
}

func extractAudio() error {
	inputFilePath, err := getInputFilePath()
	if err != nil {
		return err
	}

	outputFilePath, err := getOutputFilePath()
	if err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-acodec", "copy", outputFilePath, "-y")
	// cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("extract audio with ffmpeg: %v", err)
	}

	return nil
}

func segment() error {
	segmentService, err := getService()
	if err != nil {
		return err
	}

	switch segmentService {
	case "openai-whisper":
		if err := segmentWithOpenAIWhisper(); err != nil {
			return fmt.Errorf("segment with OpenAI Whisper: %v", err)
		}
	default:
		return fmt.Errorf("segment service `%v` not recognized", segmentService)
	}

	return nil
}

func segmentWithOpenAIWhisper() error {
	inputFilePath, err := getInputFilePath()
	if err != nil {
		return err
	}

	outputFilePath, err := getOutputFilePath()
	if err != nil {
		return err
	}

	var reqBuf bytes.Buffer
	w := multipart.NewWriter(&reqBuf)
	defer w.Close()

	file, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("open input file: %v", err)
	}
	defer file.Close()

	part, err := w.CreateFormFile("file", inputFilePath)
	if err != nil {
		return fmt.Errorf("attach input file (1): %v", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return fmt.Errorf("attach input file (2): %v", err)
	}

	writeMultipartField := func(fieldname string, value string) error {
		if err := w.WriteField(fieldname, value); err != nil {
			return fmt.Errorf("write multi-part field `%v`: %v", fieldname, err)
		}
		return nil
	}

	if err = writeMultipartField("timestamp_granularities[]", "word"); err != nil {
		return err
	}
	if err = writeMultipartField("model", "whisper-1"); err != nil {
		return err
	}
	if err = writeMultipartField("response_format", "verbose_json"); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("write req body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &reqBuf)
	if err != nil {
		return fmt.Errorf("create req: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.OpenAIApiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send req: %v", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read res body: %v", err)
	}

	os.WriteFile(outputFilePath, resBody, 0644)

	return nil
}
