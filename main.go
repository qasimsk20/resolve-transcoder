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
	"time"

	progressbar "github.com/schollz/progressbar/v3"
)

var version = "v0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Handle version flag
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("resolve-transcoder %s\n", version)
		os.Exit(0)
	}

	// Handle help flag
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printHelp()
		os.Exit(0)
	}

	inputPath := os.Args[1]

	// Check if the input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Printf("Error: Input file '%s' does not exist.\n", inputPath)
		os.Exit(1)
	}

	// Check for supported file formats
	supportedExts := []string{".mp4", ".mkv", ".avi", ".mov", ".m4v"}
	if !isSupportedFormat(inputPath, supportedExts) {
		fmt.Printf("Error: Input file '%s' format not supported.\n", inputPath)
		fmt.Printf("Supported formats: %s\n", strings.Join(supportedExts, ", "))
		os.Exit(1)
	}

	// Check if ffmpeg and ffprobe are available
	if err := checkDependencies(); err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Please install FFmpeg: https://ffmpeg.org/download.html")
		os.Exit(1)
	}

	// Construct output file path
	outputPath := generateOutputPath(inputPath)

	// Check if output file already exists
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("Warning: Output file '%s' already exists. Overwrite? (y/N): ", outputPath)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Transcoding cancelled.")
			os.Exit(0)
		}
	}

	fmt.Printf("Transcoding '%s' to '%s' for DaVinci Resolve compatibility...\n", inputPath, outputPath)

	// Get video info
	duration, fps, err := getVideoInfo(inputPath)
	if err != nil {
		fmt.Printf("Error getting video info: %v\n", err)
		os.Exit(1)
	}

	totalFrames := int(duration * fps)
	bar := progressbar.NewOptions(totalFrames,
		progressbar.OptionSetDescription("Transcoding"),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("fps"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Start transcoding
	startTime := time.Now()
	if err := transcode(inputPath, outputPath, bar); err != nil {
		fmt.Printf("Error during transcoding: %v\n", err)
		// Clean up partial output file
		os.Remove(outputPath)
		os.Exit(1)
	}

	duration_elapsed := time.Since(startTime)
	fmt.Printf("\n✓ Transcoding complete in %v!\n", duration_elapsed.Round(time.Second))
	fmt.Printf("Output: %s\n", outputPath)
	fmt.Println("You can now import the transcoded file into DaVinci Resolve.")
}

func printUsage() {
	fmt.Println("Usage: resolve-transcoder [options] <input_file_path>")
	fmt.Println("Use --help for more information.")
}

func printHelp() {
	fmt.Printf("resolve-transcoder %s\n\n", version)
	fmt.Println("A tool to transcode video files for DaVinci Resolve compatibility.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  resolve-transcoder <input_file_path>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help     Show this help message")
	fmt.Println("  -v, --version  Show version information")
	fmt.Println()
	fmt.Println("Supported input formats:")
	fmt.Println("  .mp4, .mkv, .avi, .mov, .m4v")
	fmt.Println()
	fmt.Println("Output:")
	fmt.Println("  DNxHR HQ codec in .mov container with ALAC audio")
	fmt.Println()
	fmt.Println("Requirements:")
	fmt.Println("  FFmpeg must be installed and available in PATH")
}

func isSupportedFormat(path string, supportedExts []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, supported := range supportedExts {
		if ext == supported {
			return true
		}
	}
	return false
}

func checkDependencies() error {
	// Check for ffmpeg
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH")
	}

	// Check for ffprobe
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return fmt.Errorf("ffprobe not found in PATH")
	}

	return nil
}

func generateOutputPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	fileName := filepath.Base(inputPath)
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	outputFileName := baseName + "_resolve.mov"
	return filepath.Join(dir, outputFileName)
}

func getVideoInfo(inputPath string) (duration, fps float64, err error) {
	// Get duration
	durationCmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	)
	durationOutput, err := durationCmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get duration: %w", err)
	}

	durationStr := strings.TrimSpace(string(durationOutput))
	duration, err = strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	// Get FPS
	fpsCmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=r_frame_rate",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	)
	fpsOutput, err := fpsCmd.Output()
	if err != nil {
		// Fallback to 30 FPS if we can't get the actual FPS
		return duration, 30.0, nil
	}

	fpsStr := strings.TrimSpace(string(fpsOutput))
	if strings.Contains(fpsStr, "/") {
		parts := strings.Split(fpsStr, "/")
		if len(parts) == 2 {
			num, err1 := strconv.ParseFloat(parts[0], 64)
			den, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 == nil && err2 == nil && den != 0 {
				fps = num / den
			} else {
				fps = 30.0 // Fallback
			}
		}
	} else {
		fps, err = strconv.ParseFloat(fpsStr, 64)
		if err != nil {
			fps = 30.0 // Fallback
		}
	}

	return duration, fps, nil
}

func transcode(inputPath, outputPath string, bar *progressbar.ProgressBar) error {
	// ffmpeg command with better settings
	cmd := exec.Command(
		"ffmpeg",
		"-v", "error",
		"-stats",
		"-i", inputPath,
		"-c:v", "dnxhd",
		"-profile:v", "dnxhr_hq",
		"-pix_fmt", "yuv422p",
		"-c:a", "alac",
		"-y", // Overwrite output file
		"-progress", "pipe:2",
		outputPath,
	)

	cmd.Stdout = io.Discard

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Parse progress
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		frameRegex := regexp.MustCompile(`frame=\s*(\d+)`)
		
		for scanner.Scan() {
			line := scanner.Text()
			if matches := frameRegex.FindStringSubmatch(line); len(matches) > 1 {
				if frame, err := strconv.Atoi(matches[1]); err == nil {
					bar.Set(frame)
				}
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	bar.Finish()
	return nil
}
