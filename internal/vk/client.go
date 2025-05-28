package vk

import (
	"fmt"
	"log"
	"strings"
)

// NewClient creates a new VK API client
func NewClient(appID, ownerID, token string) *Client {
	return &Client{
		AppID:   appID,
		OwnerID: ownerID,
		token:   token,
	}
}

// WallPost creates a new wall post with optional message, photos, and videos
func (c *Client) WallPost(message string, photos, videos []string, donutPaidDuration ...string) error {
	attachments := []string{}

	// Upload photos if any
	for _, photoPath := range photos {
		log.Printf("Uploading photo: %s", photoPath)
		ownerID, photoID, err := c.UploadPhoto(photoPath)
		if err != nil {
			log.Printf("Failed to upload photo %s: %v", photoPath, err)
			continue
		}
		attachments = append(attachments, fmt.Sprintf("photo%d_%d", ownerID, photoID)) // Changed from %s to %d for photoID
	}

	// Upload videos if any
	for _, videoPath := range videos {
		log.Printf("Uploading video: %s", videoPath)
		ownerID, videoID, err := c.UploadVideo(videoPath)
		if err != nil {
			log.Printf("Failed to upload video %s: %v", videoPath, err)
			continue
		}
		attachments = append(attachments, fmt.Sprintf("video%d_%d", ownerID, videoID)) // Changed from %s to %d for videoID
	}

	// Post to wall
	attachmentsStr := strings.Join(attachments, ",")
	log.Printf("Posting to wall with attachments: %s", attachmentsStr)
	return c.postToWall(message, attachmentsStr, donutPaidDuration[0])
}
