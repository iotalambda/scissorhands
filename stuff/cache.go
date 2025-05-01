package stuff

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const BaseCacheDirPath = ".scissorhandscache"

func CalculateCacheDirPath(inputFilePath string) (string, error) {
	file, err := os.Open(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("open input file: %v", err)
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("buffer input file: %v", err)
	}

	hashBytes := hash.Sum(nil)[:]
	hashStr := hex.EncodeToString(hashBytes)

	cacheDirPath := filepath.Join(BaseCacheDirPath, hashStr)

	return cacheDirPath, nil
}

func EnsureCacheDir(inputFilePath string) (string, error) {
	cacheDirPath, err := CalculateCacheDirPath(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("calculate cache dir path: %v", err)
	}

	info, err := os.Stat(cacheDirPath)
	if err != nil {
		return "", fmt.Errorf("cache dir path stat: %v", err)
	}

	if os.IsExist(err) {
		if !info.IsDir() {
			return "", fmt.Errorf("cache dir path is not a dir: %v", cacheDirPath)
		}
	} else {
		if err = os.MkdirAll(cacheDirPath, 0755); err != nil {
			return "", fmt.Errorf("make cache dir: %v", err)
		}
	}

	return cacheDirPath, nil
}
