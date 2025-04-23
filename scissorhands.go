package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {

	getArg := func(p int) (string, error) {
		if len(os.Args) <= p {
			return "", fmt.Errorf("arg at position %v not provided", p)
		}
		return os.Args[p], nil
	}

	getInputFilePath := func(p int) (string, error) {
		inputFilePath, err := getArg(p)
		if err != nil {
			return "", fmt.Errorf("error getting the input file path: %v", err)
		}
		return inputFilePath, nil
	}

	getOutputFilePath := func(p int) (string, error) {
		outputFilePath, err := getArg(p)
		if err != nil {
			return "", fmt.Errorf("error getting the output file path: %v", err)
		}
		return outputFilePath, nil
	}

	op, err := getArg(1)
	if err != nil {
		log.Fatalf("error getting the operation: %v", err)
	}

	switch op {
	case "extract-audio":
		inputFilePath, err := getInputFilePath(2)
		if err != nil {
			log.Fatalf("%v", err)
		}

		outputFilePath, err := getOutputFilePath(3)
		if err != nil {
			log.Fatalf("%v", err)
		}

		cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-acodec", "copy", outputFilePath, "-y")
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatalf("could not extract audio with ffmpeg: %v", err)
		}

	case "diarize":
		fmt.Println("TODO: diarize")
	default:
		log.Fatalf("operation `%v` not recognized", op)
	}
}
