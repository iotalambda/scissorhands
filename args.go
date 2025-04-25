package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func getArgAt(i int) (string, error) {
	if len(os.Args) <= i {
		return "", fmt.Errorf("arg at the ix %v not provided", i)
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

func getOp() (string, error) {
	op, err := getArgAt(1)
	if err != nil {
		return "", fmt.Errorf("error getting the operation: %v", err)
	}
	return op, nil
}

var inputFilePathFlags = [...]string{"-i", "--input"}

func getInputFilePath() (string, error) {
	inputFilePath, err := getArg(inputFilePathFlags[:]...)
	if err != nil {
		return "", fmt.Errorf("get input file path arg: %v", err)
	}
	return inputFilePath, nil
}

var outputFilePathFlags = [...]string{"-o", "--output"}

func getOutputFilePath() (string, error) {
	outputFilePath, err := getArg(outputFilePathFlags[:]...)
	if err != nil {
		return "", fmt.Errorf("get output file path arg: %v", err)
	}
	return outputFilePath, nil
}

var serviceFlags = [...]string{"-s", "--service"}

func getService() (string, error) {
	service, err := getArg(serviceFlags[:]...)
	if err != nil {
		return "", fmt.Errorf("get service arg: %v", err)
	}
	return service, nil
}
