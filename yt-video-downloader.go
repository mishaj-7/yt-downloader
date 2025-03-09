package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kkdai/youtube/v2"
)

func main() {
	// YouTube video URL
	videoURL := "https://youtu.be/9kSbUqv6EEo?si=cCtggzwvEwQGowZ-" // Changed to a public video for testing

	ctx := context.Background()

	// Initialize the YouTube client
	client := youtube.Client{}

	// Get video information
	video, err := client.GetVideoContext(ctx, videoURL)
	if err != nil {
		log.Fatalf("Error getting video info: %v", err)
	}

	// Print all available formats for debugging
	fmt.Println("Available formats:")
	for i, format := range video.Formats {
		fmt.Printf("Format %d: ItagNo=%d, Quality=%s, MimeType=%s\n", i, format.ItagNo, format.Quality, format.MimeType)
	}

	// Choose a format (prefer itag 22, fallback to first available video+audio format)
	var chosenFormat *youtube.Format
	for _, format := range video.Formats {
		if format.ItagNo == 22 { // 720p MP4
			chosenFormat = &format
			break
		}
	}

	// If itag 22 not found, pick the first format with both video and audio
	if chosenFormat == nil {
		for _, format := range video.Formats {
			if format.AudioChannels > 0 && format.Width > 0 { // Ensures it has video and audio
				chosenFormat = &format
				break
			}
		}
	}

	if chosenFormat == nil {
		log.Fatal("No suitable video format with both video and audio found.")
	}

	fmt.Printf("Chosen format: ItagNo=%d, Quality=%s, MimeType=%s\n", chosenFormat.ItagNo, chosenFormat.Quality, chosenFormat.MimeType)

	// Open the file to save the downloaded video
	outFile, err := os.Create("downloaded_video.mp4")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer outFile.Close()

	// Download the video
	stream, _, err := client.GetStreamContext(ctx, video, chosenFormat)
	if err != nil {
		log.Fatalf("Error getting video stream: %v", err)
	}

	_, err = outFile.ReadFrom(stream)
	if err != nil {
		log.Fatalf("Error writing video to file: %v", err)
	}

	fmt.Println("Download successful: downloaded_video.mp4")
}