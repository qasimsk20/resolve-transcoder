package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	progressbar "github.com/schollz/progressbar/v3"
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
	if !strings.HasSuffix(strings.ToLower(inputPath), ".mp4") && !strings.HasSuffix(strings.ToLower(inputPath), ".mkv") {
		fmt.Printf("Error: Input file '%s' is not an MP4 or MKV. This transcoder is designed for MP4 and MKV files.\n", inputPath)
		os.Exit(1)
	}

	// Construct output file path
	dir := filepath.Dir(inputPath)
	fileName := filepath.Base(inputPath)
	outputFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "_transcoded.mov"
	outputPath := filepath.Join(dir, outputFileName)

	fmt.Printf("Transcoding '%s' to '%s' for DaVinci Resolve compatibility...\n", inputPath, outputPath)

	// Get duration of the input video using ffprobe (part of ffmpeg suite)
	durationCmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	)
	durationOutput, err := durationCmd.Output()
	if err != nil {
		fmt.Printf("Error getting video duration: %v\n", err)
		os.Exit(1)
	}
	durationStr := strings.TrimSpace(string(durationOutput))
	durationFloat, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		fmt.Printf("Error parsing video duration: %v\n", err)
		os.Exit(1)
	}
	totalFrames := int(durationFloat * 30) // Assuming 30 FPS for progress estimation

	bar := progressbar.Default(int64(totalFrames), "Transcoding")

	// ffmpeg command and arguments
	cmd := exec.Command(
		"ffmpeg",
		"-v", "quiet",
		"-i", inputPath,
		"-c:v", "dnxhd",
		"-profile:v", "dnxhr_hq",
		"-pix_fmt", "yuv422p",
		"-c:a", "alac",
		"-progress", "pipe:2", // Output progress to stderr
		outputPath,
	)

	// Discard stdout to prevent ffmpeg's verbose output
	cmd.Stdout = io.Discard

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		os.Exit(1)
	}

	// Start ffmpeg command
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting ffmpeg: %v\n", err)
		os.Exit(1)
	}

	// Goroutine to read ffmpeg progress
	scanner := bufio.NewScanner(stderrPipe)
	go func() {
		frameRegex := regexp.MustCompile(`frame=(\d+)`)
		for scanner.Scan() {
			line := scanner.Text()
			// Only process lines that contain progress information
			if matches := frameRegex.FindStringSubmatch(line); len(matches) > 1 {
				if frame, err := strconv.Atoi(matches[1]); err == nil {
					bar.Set(frame)
				}
			}
		}
		bar.Finish()
	}()

	// Wait for ffmpeg to complete
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error during transcoding: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nTranscoding complete! You can now import the transcoded file into DaVinci Resolve.")
}
