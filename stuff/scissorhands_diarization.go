package stuff

import (
	"encoding/json"
	"fmt"
	"os"
)

type ScissorHandsDiarizationFile struct {
	DurationMilliseconds int                                         `json:"durationMilliseconds"`
	CombinedPhrases      []ScissorHandsDiarizationFileCombinedPhrase `json:"combinedPhrases"`
	Phrases              []ScissorHandsDiarizationFilePhrase         `json:"phrases"`
}

type ScissorHandsDiarizationFileCombinedPhrase struct {
	Text string `json:"text"`
}

type ScissorHandsDiarizationFilePhrase struct {
	Speaker              string                                  `json:"speaker"`
	OffsetMilliseconds   int                                     `json:"offsetMilliseconds"`
	DurationMilliseconds int                                     `json:"durationMilliseconds"`
	Text                 string                                  `json:"text"`
	Words                []ScissorHandsDiarizationFilePhraseWord `json:"words"`
	Locale               string                                  `json:"locale"`
	Confidence           float32                                 `json:"confidence"`
}

type ScissorHandsDiarizationFilePhraseWord struct {
	Text                 string `json:"text"`
	OffsetMilliseconds   int    `json:"offsetMilliseconds"`
	DurationMilliseconds int    `json:"durationMilliseconds"`
}

func ScissorHandsDiarizationFileRead(filePath string) (*ScissorHandsDiarizationFile, error) {
	bytesFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var tgt ScissorHandsDiarizationFile
	if err = json.Unmarshal(bytesFile, &tgt); err != nil {
		return nil, fmt.Errorf("unmarshal file: %v", err)
	}

	return &tgt, nil
}

func (file *ScissorHandsDiarizationFile) Write(filePath string) error {
	bytesFile, err := json.Marshal(file)
	if err != nil {
		return fmt.Errorf("marshal file: %v", err)
	}

	if err = os.WriteFile(filePath, bytesFile, 0644); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}
