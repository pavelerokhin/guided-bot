package main

import (
	"encoding/base64"
	"os"
	"strings"
)

func SaveBase64Image(base64String, outputPath string) error {
	// Decode the base64 string into bytes
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(base64String))
	if err != nil {
		return err
	}

	// Create or overwrite the file at the specified output path
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the decoded data to the file
	_, err = file.Write(decoded)
	if err != nil {
		return err
	}

	return nil
}
