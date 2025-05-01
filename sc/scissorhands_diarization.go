package sc

import (
	"encoding/json"
	"fmt"
	"os"
)

type Diarization struct {
	DurationMilliseconds int                         `json:"durationMilliseconds"`
	CombinedPhrases      []DiarizationCombinedPhrase `json:"combinedPhrases"`
	Phrases              []DiarizationPhrase         `json:"phrases"`
}

type DiarizationCombinedPhrase struct {
	Text string `json:"text"`
}

type DiarizationPhrase struct {
	Speaker              string                  `json:"speaker"`
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

func (file *Diarization) Write(path string) error {
	bytes, err := json.Marshal(file)
	if err != nil {
		return fmt.Errorf("marshal file: %v", err)
	}

	if err = os.WriteFile(path, bytes, 0644); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}
