# resolve-transcoder

`resolve-transcoder` is a simple command-line utility written in Go that helps DaVinci Resolve users on Linux transcode unsupported MP4 video files (typically H.264/H.265 video with AAC audio) into a DaVinci Resolve-compatible format (DNxHD video with ALAC audio in a MOV container).

This tool is particularly useful for users with Intel Integrated Graphics on Arch Linux (or similar distributions) who encounter issues importing common MP4 files into DaVinci Resolve Free or Studio versions due to codec limitations.

## Why use this tool?

DaVinci Resolve on Linux, especially the free version, has known limitations with certain common video and audio codecs like H.264, H.265, and AAC. This often leads to "media offline" errors or inability to import files. `resolve-transcoder` automates the process of converting these files into a format that DaVinci Resolve can readily use, saving you time and effort.

## Features

*   **Simple Command-Line Interface:** Easy to use by providing the input file path.
*   **DaVinci Resolve Compatibility:** Transcodes video to DNxHD and audio to ALAC, wrapped in a MOV container.
*   **Error Handling:** Checks for file existence and provides informative messages.
*   **Open Source:** Freely available for use, modification, and distribution.

## Prerequisites

Before using `resolve-transcoder`, ensure you have the following installed on your Arch Linux system:

*   **Go (Golang):** Required to build the program from source.
    ```bash
    sudo pacman -S go
    ```
*   **FFmpeg:** The underlying tool used for transcoding.
    ```bash
    sudo pacman -S ffmpeg
    ```

## Installation

Follow these steps to install `resolve-transcoder` on your system:

1.  **Clone the repository (or download the source code):**
    ```bash
    git clone https://github.com/YOUR_USERNAME/resolve-transcoder.git
    cd resolve-transcoder
    ```
    *(Replace `YOUR_USERNAME` with your actual GitHub username after you've uploaded the repository.)*

2.  **Build the executable:**
    ```bash
    go build -o resolve-transcoder
    ```
    This will create an executable file named `resolve-transcoder` in the current directory.

3.  **Install the executable to your system's PATH:**
    To make the program accessible from any directory, move it to a directory included in your system's `PATH` (e.g., `/usr/local/bin`).
    ```bash
    sudo mv resolve-transcoder /usr/local/bin/
    ```

4.  **Verify installation:**
    You can test if the program is correctly installed by running:
    ```bash
    resolve-transcoder
    ```
    It should output the usage message: `Usage: resolve-transcoder <input_file_path>`

## Usage

To transcode an MP4 file, simply run `resolve-transcoder` followed by the path to your input MP4 file. If the path contains spaces, enclose it in double quotes.

```bash
resolve-transcoder "/path/to/your/video with spaces.mp4"
```

**Example:**

```bash
resolve-transcoder "/home/vermillion/system 83/AQOsIDbBzoopqt-wY5GYkbm4z5KIRhuxjuPhHE6SAr5rT27ljEEiHKZmrLX8oNs-ZAIgd3yAeUcAtDglJTEz4QQT5OeHhdgjaW3hCO0.mp4"
```

The transcoded file will be saved in the same directory as the input file, with `_transcoded.mov` appended to its name.

For example, if your input was `my_video.mp4`, the output will be `my_video_transcoded.mov`.

## Contributing

Contributions are welcome! If you find a bug or have a feature request, please open an issue on the GitHub repository. If you'd like to contribute code, please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
