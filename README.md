<h1 align = center> resolve-transcoder </h1>
<div align="center">
  
![GitHub last commit](https://img.shields.io/github/last-commit/qasimsk20/supersonic.nvim?style=for-the-badge&labelColor=101418&color=%2389b4fa)
![GitHub Repo stars](https://img.shields.io/github/stars/qasimsk20/supersonic.nvim?style=for-the-badge&labelColor=101418&color=%23cba6f7)
![Repo size](https://img.shields.io/github/repo-size/qasimsk20/supersonic.nvim?style=for-the-badge&labelColor=101418&color=%23d3bfe6)
![License](https://img.shields.io/github/license/qasimsk20/supersonic.nvim?style=for-the-badge&labelColor=101418&color=%23cba6f7)

</div>

`resolve-transcoder` is a simple command-line utility written in Go that helps DaVinci Resolve users on Linux transcode unsupported MP4 video files (typically H.264/H.265 video with AAC audio) into a DaVinci Resolve-compatible format (DNxHD video with ALAC audio in a MOV container).

This tool is particularly useful for users with Intel Integrated Graphics on Arch Linux (or similar distributions) who encounter issues importing common MP4 files into DaVinci Resolve Free or Studio versions due to codec limitations.

## Why use this tool?

DaVinci Resolve on Linux, especially the free version, has known limitations with certain common video and audio codecs like H.264, H.265, and AAC. This often leads to "media offline" errors or inability to import files. `resolve-transcoder` automates the process of converting these files into a format that DaVinci Resolve can readily use, saving you time and effort.

## Features

- **Simple Command-Line Interface:** Easy to use by providing the input file path
- **DaVinci Resolve Compatibility:** Transcodes video to DNxHD and audio to ALAC, wrapped in a MOV container
- **Real-time Progress Tracking:** Visual progress bar with ETA and speed information
- **Multiple Format Support:** Works with MP4, MKV, AVI, MOV, and M4V files
- **Cross-platform:** Works on Linux, macOS, and Windows
- **Error Handling:** Checks for file existence and provides informative messages
- **Dependency Checking:** Verifies FFmpeg installation before starting

## Prerequisites

Before using `resolve-transcoder`, ensure you have the following installed on your system:

### Arch Linux
*   **Go (Golang):** Required to build the program from source.
    ```bash
    sudo pacman -S go
    ```
*   **FFmpeg:** The underlying tool used for transcoding.
    ```bash
    sudo pacman -S ffmpeg
    ```

### Other Linux Distributions
- Ubuntu/Debian: `sudo apt install ffmpeg golang-go`
- CentOS/RHEL: `sudo yum install ffmpeg golang`

### macOS
- Install via Homebrew: `brew install ffmpeg go`

### Windows
- Download FFmpeg from [ffmpeg.org](https://ffmpeg.org/download.html)
- Download Go from [golang.org](https://golang.org/dl/)

## Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Go to the [Releases page](https://github.com/qasimsk20/resolve-transcoder/releases)
2. Download the appropriate version for your platform
3. Extract the archive
4. Move the executable to your PATH:
   ```bash
   sudo mv resolve-transcoder /usr/local/bin/
   ```

### Option 2: Build from Source

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/qasimsk20/resolve-transcoder.git
    cd resolve-transcoder
    ```

2.  **Build the executable:**
    ```bash
    go build -o resolve-transcoder
    ```

3.  **Install to system PATH:**
    ```bash
    sudo mv resolve-transcoder /usr/local/bin/
    ```

4.  **Verify installation:**
    ```bash
    resolve-transcoder --version
    ```

## Usage

### Basic Usage
```bash
resolve-transcoder input_video.mp4
```

This will create `input_video_resolve.mov` in the same directory.

### Command-line Options
```bash
resolve-transcoder --help        # Show help information
resolve-transcoder --version     # Show version number
```

### Examples

```bash
# Transcode a single file
resolve-transcoder my_video.mp4

# Transcode with full path (use quotes for spaces)
resolve-transcoder "/path/to/video with spaces.mkv"

# Real example
resolve-transcoder "/home/vermillion/system 83/AQOsIDbBzoopqt-wY5GYkbm4z5KIRhuxjuPhHE6SAr5rT27ljEEiHKZmrLX8oNs-ZAIgd3yAeUcAtDglJTEz4QQT5OeHhdgjaW3hCO0.mp4"
```

The transcoded file will be saved in the same directory as the input file, with `_resolve.mov` appended to its name.

## Output Format

The transcoded files use these specifications optimized for DaVinci Resolve:

- **Video Codec:** DNxHR HQ (high quality, edit-friendly)
- **Audio Codec:** ALAC (lossless)
- **Container:** MOV (QuickTime)
- **Pixel Format:** YUV 4:2:2 (professional broadcast standard)

## Development

### Building for Multiple Platforms
```bash
make build-all    # Build for all platforms
make release      # Create release packages
```

### Running Tests
```bash
make test
```
