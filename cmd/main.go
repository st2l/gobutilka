package main

import (
	"log"
	"time"

	"github.com/st2l/vk_butilka/internal/config"
	"github.com/st2l/vk_butilka/internal/content"
	"github.com/st2l/vk_butilka/internal/vk"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create VK client
	client := vk.NewClient(cfg.AppID, cfg.OwnerID, cfg.VKAccountToken)

	// Initialize content manager using config values
	manager := content.NewContentManager(cfg.ContentDir, cfg.DonutFrequency, cfg.DonutDuration, cfg.ContentPerPost)

	// Use configured post interval
	interval := cfg.PostInterval

	log.Printf("Starting automatic posting service with %d hour intervals", interval/time.Hour)
	log.Printf("Donut posts will appear every %d regular posts", cfg.DonutFrequency)
	log.Printf("Monitoring content directory: %s", cfg.ContentDir)

	// Run first post immediately
	makePost(client, manager)

	// Then set up a ticker for periodic posts
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		makePost(client, manager)
	}
}

func makePost(client *vk.Client, manager *content.ContentManager) {
	// Scan for content
	photos, videos, err := manager.ScanContent()
	if err != nil {
		log.Printf("Failed to scan content: %v", err)
		return
	}

	if len(photos) == 0 && len(videos) == 0 {
		log.Printf("No media content found. Skipping this post cycle.")
		return
	}

	// Check if we should use donut for this post
	useDonut, donutDuration := manager.ShouldUseDonut()

	// TODO: Implement logic of creating a message with AI
	message := ""
	if useDonut {
		message = ""
		log.Println("Creating a donut post...")
	}

	// Post to wall
	var postErr error
	if useDonut {
		postErr = client.WallPost(message, photos, videos, donutDuration)
	} else {
		postErr = client.WallPost(message, photos, videos, "")
	}

	if postErr != nil {
		log.Printf("Failed to post to wall: %v", postErr)
		return
	}

	log.Println("Post successful!")

	// Delete used content
	allContent := append(photos, videos...)
	manager.DeleteUsedContent(allContent)
}
