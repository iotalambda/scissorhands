package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"scissorhands/stuff"

	"github.com/spf13/cobra"
)

var inferDialogueCmd = &cobra.Command{
	Use:   "infer-dialogue",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return inferDialogue()
	},
}

func inferDialogue() error {
	cacheDirPath, err := stuff.EnsureCacheDir(input)
	if err != nil {
		return fmt.Errorf("ensure cache dir: %v", err)
	}

	audioFilePath := filepath.Join(cacheDirPath, "audio.aac")
	_, err = os.Stat(audioFilePath)
	if os.IsNotExist(err) {
		if err = stuff.FfmpegExtractAudio(input, audioFilePath); err != nil {
			return fmt.Errorf("ffmpeg extract audio: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("audio file stat: %v", err)
	}

	azureSpeechDiarizationFilePath := filepath.Join(cacheDirPath, "azure_speech_diarization.json")
	_, err = os.Stat(azureSpeechDiarizationFilePath)
	if os.IsNotExist(err) {
		if err = stuff.AzureSpeechDiarize(audioFilePath, azureSpeechDiarizationFilePath, 5); err != nil {
			return fmt.Errorf("azure speech diarize: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("azure speech diarization file stat: %v", err)
	}

	scissorhandsDiarizationFilePath := filepath.Join(cacheDirPath, "scissorhands_diarization.json")
	_, err = os.Stat(scissorhandsDiarizationFilePath)
	if os.IsNotExist(err) {
		azureSpeechDiarizationFile, err := stuff.AzureSpeechDiarizationFileRead(azureSpeechDiarizationFilePath)
		if err != nil {
			return fmt.Errorf("read azure speech diarization file: %v", err)
		}

		scissorhandsDiarizationFile, err := azureSpeechDiarizationFile.MapToScissorhandsDiarization()
		if err != nil {
			return fmt.Errorf("map azure speech to scissorhands diarization file: %v", err)
		}

		for pIx := range scissorhandsDiarizationFile.Phrases {
			p := &scissorhandsDiarizationFile.Phrases[pIx]
			p.Speaker = ""
		}

		if err = scissorhandsDiarizationFile.Write(scissorhandsDiarizationFilePath); err != nil {
			return fmt.Errorf("write scissorhands diarization file: %v", err)
		}

	} else if err != nil {
		return fmt.Errorf("scissorhands diarization file stat: %v", err)
	}

	return nil
}

func init() {
	inferDialogueCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	inferDialogueCmd.MarkFlagRequired("input")
}
