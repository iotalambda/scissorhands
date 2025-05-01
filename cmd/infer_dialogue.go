package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"scissorhands/azspeech"
	"scissorhands/cache"
	"scissorhands/ffmpeg"
	"scissorhands/sc"

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
	cacheDirPath, err := cache.EnsureCacheDir(input)
	if err != nil {
		return fmt.Errorf("ensure cache dir: %v", err)
	}

	audioPath := filepath.Join(cacheDirPath, "audio.aac")
	_, err = os.Stat(audioPath)
	if os.IsNotExist(err) {
		if err = ffmpeg.ExtractAudio(input, audioPath); err != nil {
			return fmt.Errorf("ffmpeg extract audio: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("audio file stat: %v", err)
	}

	azureSpeechDiarizationPath := filepath.Join(cacheDirPath, "azure_speech_diarization.json")
	_, err = os.Stat(azureSpeechDiarizationPath)
	var azureSpeechDiarization *azspeech.Diarization
	if os.IsNotExist(err) {
		azureSpeechDiarization, err = azspeech.Diarize(audioPath, 5)
		if err != nil {
			return fmt.Errorf("azure speech diarize: %v", err)
		}
		azureSpeechDiarization.Write(azureSpeechDiarizationPath)
	} else if err != nil {
		return fmt.Errorf("azure speech diarization file stat: %v", err)
	} else {
		azureSpeechDiarization, err = azspeech.ReadDiarization(azureSpeechDiarizationPath)
		if err != nil {
			return fmt.Errorf("azure speech read diarization: %v", err)
		}
	}

	scissorhandsDiarizationPath := filepath.Join(cacheDirPath, "scissorhands_diarization.json")
	_, err = os.Stat(scissorhandsDiarizationPath)
	var scissorhandsDiarization *sc.Diarization
	if os.IsNotExist(err) {
		scissorhandsDiarization, err = azureSpeechDiarization.MapToScissorhands()
		if err != nil {
			return fmt.Errorf("map azure speech to scissorhands diarization file: %v", err)
		}

		// Clear Azure Speech speaker information. It's unreliable.
		for pIx := range scissorhandsDiarization.Phrases {
			p := &scissorhandsDiarization.Phrases[pIx]
			p.Speaker = ""
		}

		if err = scissorhandsDiarization.Write(scissorhandsDiarizationPath); err != nil {
			return fmt.Errorf("write scissorhands diarization file: %v", err)
		}

	} else if err != nil {
		return fmt.Errorf("scissorhands diarization file stat: %v", err)
	} else {
		scissorhandsDiarization, err = sc.ReadDiarization(scissorhandsDiarizationPath)
		if err != nil {
			return fmt.Errorf("scissorhands read diarization: %v", err)
		}
	}

	durationMs, err := ffmpeg.DurationMs(input)
	if err != nil {
		return fmt.Errorf("ffprobe duration ms: %v", err)
	}

	nScreenshots := 10
	htmlStr := "<html>"
	for i := range nScreenshots {
		timeScreenshotMs := (durationMs / nScreenshots) * i
		timeScreenshotSs := ffmpeg.MsToSeek(timeScreenshotMs)
		screenshot, err := ffmpeg.Screenshot(timeScreenshotSs, input)
		if err != nil {
			return fmt.Errorf("ffmpeg screenshot %v: %v", i, err)
		}
		htmlStr += fmt.Sprintf("<br/><img src=\"%v\" />", screenshot)
	}
	htmlStr += "</html>"
	//os.WriteFile(filepath.Join(".temp", "output", "site.html"), []byte(htmlStr), 0644)

	return nil
}

func init() {
	inferDialogueCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	inferDialogueCmd.MarkFlagRequired("input")
}
