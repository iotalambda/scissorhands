package azspeech

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"scissorhands/config"
	"scissorhands/sc"
)

func Diarize(srcPath string, maxSpeakers int) (*Diarization, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	defer w.Close()

	file, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("open source file: %v", err)
	}
	defer file.Close()

	part, err := w.CreateFormFile("audio", srcPath)
	if err != nil {
		return nil, fmt.Errorf("attach source file (1): %v", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("attach source file (2): %v", err)
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
		return nil, fmt.Errorf("marshal definition: %v", err)
	}

	if err = writeMultipartField("definition", string(definition)); err != nil {
		return nil, err
	}

	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("write req body: %v", err)
	}

	// API Docs: https://learn.microsoft.com/en-us/azure/ai-services/speech-service/fast-transcription-create?tabs=diarization-on
	req, err := http.NewRequest("POST", "https://northeurope.api.cognitive.microsoft.com/speechtotext/transcriptions:transcribe?api-version=2024-11-15", &buf)
	if err != nil {
		return nil, fmt.Errorf("create req: %v", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", config.Global.AzureSpeechServiceKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send req: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read res body: %v", err)
	}

	var tgt Diarization
	if err = json.Unmarshal(body, &tgt); err != nil {
		return nil, fmt.Errorf("unmarshal body: %v", err)
	}

	return &tgt, nil
}

func (d *Diarization) Write(path string) error {
	bytes, err := json.Marshal(&d)
	if err != nil {
		return fmt.Errorf("marshal diarization: %v", err)
	}
	if err = os.WriteFile(path, bytes, 0644); err != nil {
		return fmt.Errorf("write file: %v", err)
	}
	return nil
}

type Diarization struct {
	DurationMilliseconds int                         `json:"durationMilliseconds"`
	CombinedPhrases      []DiarizationCombinedPhrase `json:"combinedPhrases"`
	Phrases              []DiarizationPhrase         `json:"phrases"`
}

type DiarizationCombinedPhrase struct {
	Text string `json:"text"`
}

type DiarizationPhrase struct {
	Speaker              int                     `json:"speaker"`
	OffsetMilliseconds   int                     `json:"offsetMilliseconds"`
	DurationMilliseconds int                     `json:"durationMilliseconds"`
	Text                 string                  `json:"text"`
	Words                []DiarizationPhraseWord `json:"words"`
	Locale               string                  `json:"locale"`
	Confidence           float32                 `json:"confidence"`
}

type DiarizationPhraseWord struct {
	Text                 string `json:"text"`
	OffsetMilliseconds   int    `json:"offsetMilliseconds"`
	DurationMilliseconds int    `json:"durationMilliseconds"`
}

func ReadDiarization(path string) (*Diarization, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var tgt Diarization
	if err = json.Unmarshal(bytes, &tgt); err != nil {
		return nil, fmt.Errorf("unmarshal file: %v", err)
	}

	return &tgt, nil
}

func (src *Diarization) MapToScissorhands() (*sc.Diarization, error) {
	tgt := sc.Diarization{
		DurationMilliseconds: src.DurationMilliseconds,
	}
	for _, combinedPhraseSrc := range src.CombinedPhrases {
		combinedPhraseTgt := sc.DiarizationCombinedPhrase(combinedPhraseSrc)
		tgt.CombinedPhrases = append(tgt.CombinedPhrases, combinedPhraseTgt)
	}
	for _, phraseSrc := range src.Phrases {
		phraseTgt := sc.DiarizationPhrase{
			Speaker:              fmt.Sprint(phraseSrc.Speaker),
			DurationMilliseconds: phraseSrc.DurationMilliseconds,
			OffsetMilliseconds:   phraseSrc.OffsetMilliseconds,
			Text:                 phraseSrc.Text,
			Locale:               phraseSrc.Locale,
			Confidence:           phraseSrc.Confidence,
		}
		for _, wordSrc := range phraseSrc.Words {
			wordTgt := sc.DiarizationPhraseWord(wordSrc)
			phraseTgt.Words = append(phraseTgt.Words, wordTgt)
		}
		tgt.Phrases = append(tgt.Phrases, phraseTgt)
	}

	return &tgt, nil
}
