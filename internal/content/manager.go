package content

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

// SupportedExtensions defines what file types we'll consider
var SupportedExtensions = map[string]string{
	// Images
	".jpg":  "photo",
	".jpeg": "photo",
	".png":  "photo",
	".gif":  "photo",
	// Videos
	".mp4": "video",
	".mov": "video",
	".avi": "video",
}

// ContentManager handles content discovery and tracking posts
type ContentManager struct {
	contentDir     string
	postsCounter   int32
	donutFrequency int32
	donutDuration  string
	contentPerPost int
}

// NewContentManager creates a new content manager
func NewContentManager(contentDir string, donutFrequency int, donutDuration string, contentPerPost int) *ContentManager {
	return &ContentManager{
		contentDir:     contentDir,
		donutFrequency: int32(donutFrequency),
		donutDuration:  donutDuration,
		contentPerPost: contentPerPost,
	}
}

// ScanContent scans the content directory and returns limited photos and videos
func (cm *ContentManager) ScanContent() (photos []string, videos []string, err error) {
	entries, err := os.ReadDir(cm.contentDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read content directory: %w", err)
	}

	if len(entries) == 0 {
		return nil, nil, errors.New("no content found in directory")
	}

	// Gather all available content first
	var allPhotos []string
	var allVideos []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))

		if mediaType, supported := SupportedExtensions[ext]; supported {
			fullPath := filepath.Join(cm.contentDir, fileName)

			switch mediaType {
			case "photo":
				allPhotos = append(allPhotos, fullPath)
			case "video":
				allVideos = append(allVideos, fullPath)
			}
		}
	}

	// Randomize the order to avoid always picking the same files
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allPhotos), func(i, j int) {
		allPhotos[i], allPhotos[j] = allPhotos[j], allPhotos[i]
	})
	rand.Shuffle(len(allVideos), func(i, j int) {
		allVideos[i], allVideos[j] = allVideos[j], allVideos[i]
	})

	// Calculate how many of each type to include based on available content
	totalItems := cm.contentPerPost
	maxPhotos := min(len(allPhotos), totalItems)

	// Take photos up to max
	photos = allPhotos[:maxPhotos]

	// Use remaining slots for videos
	remainingSlots := totalItems - maxPhotos
	maxVideos := min(len(allVideos), remainingSlots)
	if maxVideos > 0 {
		videos = allVideos[:maxVideos]
	}

	return photos, videos, nil
}

// Helper function for the min of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DeleteUsedContent deletes files after they've been posted
func (cm *ContentManager) DeleteUsedContent(files []string) []error {
	var errors []error

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			errors = append(errors, fmt.Errorf("failed to delete %s: %w", file, err))
		}
	}

	return errors
}

// ShouldUseDonut returns true if the current post should be a donut post
func (cm *ContentManager) ShouldUseDonut() (bool, string) {
	counter := atomic.AddInt32(&cm.postsCounter, 1)

	if counter%cm.donutFrequency == 0 {
		return true, cm.donutDuration
	}

	return false, ""
}
