package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func getArgAt(i int) (string, error) {
	if len(os.Args) <= i {
		return "", fmt.Errorf("arg at ix %v not provided", i)
	}
	return os.Args[i], nil
}

func getArg(flags ...string) (string, error) {
	for i, arg := range os.Args {
		if slices.Contains(flags, arg) {
			return getArgAt(i + 1)
		}
	}
	return "", fmt.Errorf("arg with the flag %v not provided", strings.Join(flags, "/"))
}

func getOpOrErr() string {
	op, err := getArgAt(1)
	if err != nil {
		log.Fatalf("error getting the operation: %v", err)
	}
	return op
}

var inputFilePathFlags = [...]string{"-i", "--input"}

func getInputFilePathOrErr() string {
	inputFilePath, err := getArg(inputFilePathFlags[:]...)
	if err != nil {
		log.Fatalf("error getting the input file path: %v", err)
	}
	return inputFilePath
}

var outputFilePathFlags = [...]string{"-o", "--output"}

func getOutputFilePathOrErr() string {
	outputFilePath, err := getArg(outputFilePathFlags[:]...)
	if err != nil {
		log.Fatalf("error getting the output file path: %v", err)
	}
	return outputFilePath
}

var diarizeServiceFlags = [...]string{"-s", "--service"}

func getDiarizeServiceOrErr() string {
	diarizeService, err := getArg(diarizeServiceFlags[:]...)
	if err != nil {
		log.Fatalf("error getting the diarize service: %v", err)
	}
	return diarizeService
}
