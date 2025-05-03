package sc

import (
	"encoding/json"
	"fmt"
	"os"
)

type Diarization struct {
	DurationMilliseconds int                 `json:"durationMilliseconds"`
	Phrases              []DiarizationPhrase `json:"phrases"`
}

type DiarizationCombinedPhrase struct {
	Text string `json:"text"`
}

type DiarizationPhrase struct {
	Speaker              string  `json:"speaker"`
	OffsetMilliseconds   int     `json:"offsetMilliseconds"`
	DurationMilliseconds int     `json:"durationMilliseconds"`
	Text                 string  `json:"text"`
	Locale               string  `json:"locale"`
	Confidence           float32 `json:"confidence"`
}

func ReadDiarization(path string) (*Diarization, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read diarization: %v", err)
	}

	var d Diarization
	if err = json.Unmarshal(bytes, &d); err != nil {
		return nil, fmt.Errorf("unmarshal diarization: %v", err)
	}

	return &d, nil
}

func (d *Diarization) Write(path string) error {
	bytes, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("marshal diarization: %v", err)
	}

	if err = os.WriteFile(path, bytes, 0644); err != nil {
		return fmt.Errorf("write diarization: %v", err)
	}

	return nil
}

func (d *Diarization) String() (string, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		return "", fmt.Errorf("marshal diarization: %v", err)
	}
	return string(bytes), nil
}
