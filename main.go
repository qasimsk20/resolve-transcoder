package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: resolve-transcoder <input_file_path>")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	// Check if the input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Printf("Error: Input file '%s' does not exist.\n", inputPath)
		os.Exit(1)
	}

	// Ensure it's an MP4 file (basic check)
	if !strings.HasSuffix(strings.ToLower(inputPath), ".mp4") {
		fmt.Printf("Error: Input file '%s' is not an MP4. This transcoder is designed for MP4 files.\n", inputPath)
		os.Exit(1)
	}

	// Construct output file path
	dir := filepath.Dir(inputPath)
	fileName := filepath.Base(inputPath)
	outputFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "_transcoded.mov"
	outputPath := filepath.Join(dir, outputFileName)

	fmt.Printf("Transcoding '%s' to '%s' for DaVinci Resolve compatibility...\n", inputPath, outputPath)

	// ffmpeg command and arguments
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-c:v", "dnxhd",
		"-profile:v", "dnxhr_hq",
		"-pix_fmt", "yuv422p",
		"-c:a", "alac",
		outputPath,
	)

	// Capture stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error during transcoding: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Transcoding complete! You can now import the transcoded file into DaVinci Resolve.")
}
