package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"scissorhands/stuff"

	"github.com/spf13/cobra"
)

var azureSpeechDiarizeCmd = &cobra.Command{
	Use:   "azure-speech-diarize",
	Short: "Diarize an audio file by speakers",
	RunE: func(cmd *cobra.Command, args []string) error {
		return azureSpeechDiarize()
	},
}

func azureSpeechDiarize() error {
	var reqBuf bytes.Buffer
	w := multipart.NewWriter(&reqBuf)
	defer w.Close()

	file, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("open input file: %v", err)
	}
	defer file.Close()

	part, err := w.CreateFormFile("audio", input)
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

	definition, err := json.Marshal(map[string]any{
		"locales": []string{"en-US"},
		"diarization": map[string]any{
			"maxSpeakers": maxSpeakers,
			"enabled":     true,
		},
	})
	if err != nil {
		return fmt.Errorf("marshal definition: %v", err)
	}

	if err = writeMultipartField("definition", string(definition)); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("write req body: %v", err)
	}

	// API Docs: https://learn.microsoft.com/en-us/azure/ai-services/speech-service/fast-transcription-create?tabs=diarization-on
	req, err := http.NewRequest("POST", "https://northeurope.api.cognitive.microsoft.com/speechtotext/transcriptions:transcribe?api-version=2024-11-15", &reqBuf)
	if err != nil {
		return fmt.Errorf("create req: %v", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", stuff.Global.AzureSpeechServiceKey)
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
	rootCmd.AddCommand(azureSpeechDiarizeCmd)

	azureSpeechDiarizeCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	azureSpeechDiarizeCmd.MarkFlagRequired("input")

	azureSpeechDiarizeCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path.")
	azureSpeechDiarizeCmd.MarkFlagRequired("output")

	azureSpeechDiarizeCmd.Flags().IntVarP(&maxSpeakers, "max-speakers", "m", 5, "Maximum number of possible speakers to evaluate.")
}
