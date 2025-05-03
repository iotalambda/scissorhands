package sc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const BaseCacheDirPath = ".scissorhandscache"

func CalculateCacheDirPath(srcPath string) (string, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("open source file: %v", err)
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("buffer source file: %v", err)
	}

	hashBytes := hash.Sum(nil)[:]
	hashStr := hex.EncodeToString(hashBytes)

	cacheDirPath := filepath.Join(BaseCacheDirPath, hashStr)

	return cacheDirPath, nil
}

func EnsureCacheDir(pathInput string) (string, error) {
	cacheDirPath, err := CalculateCacheDirPath(pathInput)
	if err != nil {
		return "", fmt.Errorf("calculate cache dir path: %v", err)
	}

	info, err := os.Stat(cacheDirPath)

	if os.IsNotExist(err) {
		if err = os.MkdirAll(cacheDirPath, 0755); err != nil {
			return "", fmt.Errorf("make cache dir: %v", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("cache dir path stat: %v", err)
	} else {
		if !info.IsDir() {
			return "", fmt.Errorf("cache dir path is not a dir: %v", cacheDirPath)
		}
	}

	return cacheDirPath, nil
}

func EnsureCached(cachedFilePath string, create func() error) error {
	_, err := os.Stat(cachedFilePath)
	if os.IsNotExist(err) {
		if err = create(); err != nil {
			return fmt.Errorf("ensure cached create: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("ensure cached stat: %v", err)
	}
	return nil
}
