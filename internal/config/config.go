package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	// VK API settings
	VKAccountToken string
	AppID          string
	OwnerID        string

	// Content manager settings
	ContentDir     string        // Directory to scan for content
	DonutFrequency int           // How often to make donut posts (every N posts)
	PostInterval   time.Duration // Time between posts
	DonutDuration  string        // Duration of donut posts in days (-1 for unlimited)
	ContentPerPost int           // Maximum number of content items per post
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Required VK API settings
	token := os.Getenv("VK_ACCOUNT_TOKEN")
	if token == "" {
		return nil, errors.New("VK_ACCOUNT_TOKEN environment variable is not set")
	}

	appID := os.Getenv("appID")
	if appID == "" {
		return nil, errors.New("appID environment variable is not set")
	}

	ownerID := os.Getenv("ownerID")
	if ownerID == "" {
		return nil, errors.New("ownerID environment variable is not set")
	}

	// Content manager settings with defaults
	contentDir := os.Getenv("CONTENT_DIR")
	if contentDir == "" {
		contentDir = "./content"
	}

	donutFrequency := 5 // Default value
	if dfStr := os.Getenv("DONUT_FREQUENCY"); dfStr != "" {
		if df, err := strconv.Atoi(dfStr); err == nil && df > 0 {
			donutFrequency = df
		}
	}

	postInterval := 3 * time.Hour // Default value: 3 hours
	if piStr := os.Getenv("POST_INTERVAL_HOURS"); piStr != "" {
		if hours, err := strconv.Atoi(piStr); err == nil && hours > 0 {
			postInterval = time.Duration(hours) * time.Hour
		}
	}

	donutDuration := "-1" // Default value: unlimited
	if dd := os.Getenv("DONUT_DURATION"); dd != "" {
		donutDuration = dd
	}

	// Content per post limit
	contentPerPost := 5 // Default value: 5 items
	if cppStr := os.Getenv("CONTENT_PER_POST"); cppStr != "" {
		if cpp, err := strconv.Atoi(cppStr); err == nil && cpp > 0 {
			contentPerPost = cpp
		}
	}

	return &Config{
		VKAccountToken: token,
		AppID:          appID,
		OwnerID:        ownerID,
		ContentDir:     contentDir,
		DonutFrequency: donutFrequency,
		PostInterval:   postInterval,
		DonutDuration:  donutDuration,
		ContentPerPost: contentPerPost,
	}, nil
}
