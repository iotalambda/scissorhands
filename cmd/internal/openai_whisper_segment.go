package internal

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"scissorhands/sc"

	"github.com/spf13/cobra"
)

var OpenAIWhisperSegmentCmd = &cobra.Command{
	Use:   "openai-whisper-segment",
	Short: "Segment an audio file by words",
	RunE: func(cmd *cobra.Command, args []string) error {
		return openAIWhisperSegment()
	},
}

func openAIWhisperSegment() error {
	var reqBuf bytes.Buffer
	w := multipart.NewWriter(&reqBuf)
	defer w.Close()

	file, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("open input file: %v", err)
	}
	defer file.Close()

	part, err := w.CreateFormFile("file", input)
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

	// API Docs: https://platform.openai.com/docs/guides/speech-to-text?lang=curl
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &reqBuf)
	if err != nil {
		return fmt.Errorf("create req: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+sc.GlobalConfig.OpenAIApiKey)
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

	os.WriteFile(output, resBody, 0644)

	return nil
}

func init() {
	OpenAIWhisperSegmentCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	OpenAIWhisperSegmentCmd.MarkFlagRequired("input")

	OpenAIWhisperSegmentCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path.")
	OpenAIWhisperSegmentCmd.MarkFlagRequired("output")
}
