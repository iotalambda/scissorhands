package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	op := getOpOrErr()
	switch op {
	case "extract-audio":
		extractAudio()
	case "diarize":
		diarize()
	default:
		log.Fatalf("operation `%v` not recognized", op)
	}
}

func extractAudio() {
	inputFilePath := getInputFilePathOrErr()
	outputFilePath := getOutputFilePathOrErr()
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-acodec", "copy", outputFilePath, "-y")
	// cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not extract audio with ffmpeg: %v", err)
	}
}

func diarize() {
	diarizeService := getDiarizeServiceOrErr()
	switch diarizeService {
	case "openai-whisper":
		diarizeWithOpenAIWhisper()
	case "azure-speech":
		log.Fatalf("TODO")
	default:
		log.Fatalf("diarize service `%v` not recognized", diarizeService)
	}
}

func diarizeWithOpenAIWhisper() {
	inputFilePath := getInputFilePathOrErr()
	apiKey := "TODO"

	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open the input file: %v", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", inputFilePath)
	if err != nil {
		log.Fatalf("could not attach the input file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("could not attach the input file: %v", err)
	}

	writer.WriteField("timestamp_granularities[]", "word")
	writer.WriteField("model", "whisper-1")
	writer.WriteField("response_format", "verbose_json")

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &buf)

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "multipart/form-data")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("could not invoke the request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("could not read response body: %v", err)
	}
	fmt.Println("RESPONSE=", string(body))
}
