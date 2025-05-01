package stuff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func AzureSpeechDiarize(inputFilePath string, outputFilePath string, maxSpeakers int) error {
	var reqBuf bytes.Buffer
	w := multipart.NewWriter(&reqBuf)
	defer w.Close()

	file, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("open input file: %v", err)
	}
	defer file.Close()

	part, err := w.CreateFormFile("audio", inputFilePath)
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

	req.Header.Set("Ocp-Apim-Subscription-Key", GlobalConfig.AzureSpeechServiceKey)
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

type AzureSpeechDiarizationFile struct {
	DurationMilliseconds int                                        `json:"durationMilliseconds"`
	CombinedPhrases      []AzureSpeechDiarizationFileCombinedPhrase `json:"combinedPhrases"`
	Phrases              []AzureSpeechDiarizationFilePhrase         `json:"phrases"`
}

type AzureSpeechDiarizationFileCombinedPhrase struct {
	Text string `json:"text"`
}

type AzureSpeechDiarizationFilePhrase struct {
	Speaker              int                                    `json:"speaker"`
	OffsetMilliseconds   int                                    `json:"offsetMilliseconds"`
	DurationMilliseconds int                                    `json:"durationMilliseconds"`
	Text                 string                                 `json:"text"`
	Words                []AzureSpeechDiarizationFilePhraseWord `json:"words"`
	Locale               string                                 `json:"locale"`
	Confidence           float32                                `json:"confidence"`
}

type AzureSpeechDiarizationFilePhraseWord struct {
	Text                 string `json:"text"`
	OffsetMilliseconds   int    `json:"offsetMilliseconds"`
	DurationMilliseconds int    `json:"durationMilliseconds"`
}

func AzureSpeechDiarizationFileRead(filePath string) (*AzureSpeechDiarizationFile, error) {
	bytesFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var tgt AzureSpeechDiarizationFile
	if err = json.Unmarshal(bytesFile, &tgt); err != nil {
		return nil, fmt.Errorf("unmarshal file: %v", err)
	}

	return &tgt, nil
}

func (src *AzureSpeechDiarizationFile) MapToScissorhandsDiarization() (*ScissorHandsDiarizationFile, error) {
	tgt := ScissorHandsDiarizationFile{
		DurationMilliseconds: src.DurationMilliseconds,
	}
	for _, combinedPhraseSrc := range src.CombinedPhrases {
		combinedPhraseTgt := ScissorHandsDiarizationFileCombinedPhrase(combinedPhraseSrc)
		tgt.CombinedPhrases = append(tgt.CombinedPhrases, combinedPhraseTgt)
	}
	for _, phraseSrc := range src.Phrases {
		phraseTgt := ScissorHandsDiarizationFilePhrase{
			Speaker:              fmt.Sprint(phraseSrc.Speaker),
			DurationMilliseconds: phraseSrc.DurationMilliseconds,
			OffsetMilliseconds:   phraseSrc.OffsetMilliseconds,
			Text:                 phraseSrc.Text,
			Locale:               phraseSrc.Locale,
			Confidence:           phraseSrc.Confidence,
		}
		for _, wordSrc := range phraseSrc.Words {
			wordTgt := ScissorHandsDiarizationFilePhraseWord(wordSrc)
			phraseTgt.Words = append(phraseTgt.Words, wordTgt)
		}
		tgt.Phrases = append(tgt.Phrases, phraseTgt)
	}

	return &tgt, nil
}
