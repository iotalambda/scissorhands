package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"scissorhands/azspeech"
	"scissorhands/ffmpeg"
	"scissorhands/sc"
	"strings"

	"github.com/openai/openai-go"
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
	ctx := context.Background()
	unknownSpeaker := "<UNKNOWN>"

	cacheDirPath, err := sc.EnsureCacheDir(input)
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

		// Mark Azure Speech speaker information as unknown. The results are unreliable.
		for pIx := range scissorhandsDiarization.Phrases {
			p := &scissorhandsDiarization.Phrases[pIx]
			p.Speaker = unknownSpeaker
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

	// TODO: detect freeze frame segments -> fill narrator parts

	durationMs, err := ffmpeg.DurationMs(input)
	if err != nil {
		return fmt.Errorf("ffprobe duration ms: %v", err)
	}

	nScreenshots := 10

	// Overall description
	{
		client := sc.NewOpenAIClient()
		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are an assistant tasked with producing a concise description of a video. Your description will be used by another AI agent responsible for video editing and diarization.  
For each video, you will receive a set of screenshots, the transcription, and the video's duration.

Your response must follow this XML structure:

<VideoDescription>
    <Summary>Brief summary of the overall content</Summary>
    <Individuals>
        <Individual>
            <Name>Description or name of the individual</Name>
            <Role>Role or function (e.g., narrator, speaker, host, interviewee)</Role>
        </Individual>
        <!-- Repeat <Individual> for each individual -->
    </Individuals>
</VideoDescription>

Clearly identify all individuals present or speaking. Prioritize listing everyone who speaks or is likely to speak to assist the diarization process.  
If any part of the transcription includes dialogue phrased in the passive voice or from a third-person perspective (e.g., "Next, the procedure is performed"), and no speaker is explicitly identified, treat it as narration. Assume such narration may be provided by a narrator not shown in the video.  
Be brief but complete. Do not include unrelated commentary.`),
		}

		scissorhandsDiarizationString, err := scissorhandsDiarization.String()
		if err != nil {
			return fmt.Errorf("scissorhands diarization string: %v", err)
		}

		userMessageParts := []openai.ChatCompletionContentPartUnionParam{
			openai.TextContentPart(fmt.Sprintf("The video duration is %v ms.", durationMs)),
			openai.TextContentPart(fmt.Sprintf("Here's the video diarization. Unknown speakers are marked with %v: %v", unknownSpeaker, scissorhandsDiarizationString)),
		}

		var htmlB strings.Builder
		htmlB.WriteString("<html>")
		for i := range nScreenshots {
			timeScreenshotMs := (durationMs / nScreenshots) * i
			timeScreenshotSs := ffmpeg.MsToSeek(timeScreenshotMs)
			url, err := ffmpeg.Screenshot(timeScreenshotSs, input)
			if err != nil {
				return fmt.Errorf("ffmpeg screenshot %v: %v", i, err)
			}
			userMessageParts = append(
				userMessageParts,
				openai.TextContentPart(fmt.Sprintf("Here's screenshot %v:", i+1)),
				openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL: url,
				}))
			htmlB.WriteString(fmt.Sprintf("<br/><img src=\"%v\" />", url))
		}
		htmlB.WriteString("</html>")
		//os.WriteFile(filepath.Join(".temp", "output", "screenshots.html"), []byte(htmlB.String()), 0644)
		messages = append(messages, openai.UserMessage(userMessageParts))

		resp, err := client.Chat.Completions.New(ctx, sc.NewOpenAIChatCompletionNewParams(messages))
		if err != nil {
			return fmt.Errorf("overall description chat completion: %v", err)
		}

		fmt.Printf("LLM> %v\n", resp.Choices[0].Message.Content)
	}

	return nil
}

func init() {
	inferDialogueCmd.Flags().StringVarP(&input, "input", "i", "", "Input file path.")
	inferDialogueCmd.MarkFlagRequired("input")
}
