package vk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// UploadPhoto uploads a photo to VK and returns owner_id and photo_id
func (c *Client) UploadPhoto(photoPath string) (int, int, error) {
	// Step 1: Get wall upload server
	_, uploadURL, _, err := c.getWallUploadServer()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get upload server: %w", err)
	}

	log.Printf("Got upload URL for photo: %s", uploadURL)

	// Step 2: Upload file to server
	server, photo, hash, err := c.uploadFileToServer(uploadURL, photoPath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to upload file to server: %w", err)
	}

	log.Printf("Uploaded photo to server: server=%d", server)

	// Step 3: Save photo
	ownerID, photoID, err := c.saveWallPhoto(photo, server, hash)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to save wall photo: %w", err)
	}

	log.Printf("Saved wall photo: owner=%d, id=%d", ownerID, photoID)

	return ownerID, photoID, nil
}

// UploadVideo uploads a video to VK and returns owner_id and video_id
func (c *Client) UploadVideo(videoPath string) (int, int, error) {
	// Step 1: Get video upload URL
	_, uploadURL, err := c.videoSave()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get video upload URL: %w", err)
	}

	log.Printf("Got video upload URL: %s", uploadURL)

	// Step 2: Send video to server
	vOwnerID, videoID, err := c.uploadVideoToServer(uploadURL, videoPath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to upload video to server: %w", err)
	}

	log.Printf("Uploaded video to server: owner=%d, id=%d", vOwnerID, videoID)

	return vOwnerID, videoID, nil
}

// getWallUploadServer gets the server URL for photo uploads
func (c *Client) getWallUploadServer() (int, string, int, error) {
	url := "https://api.vk.com/method/photos.getWallUploadServer"

	// Remove the minus sign from owner ID if present
	groupID := strings.Replace(c.OwnerID, "-", "", 1)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", 0, err
	}

	q := req.URL.Query()
	q.Add("group_id", groupID)
	q.Add("access_token", c.token)
	q.Add("v", "5.199")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", 0, err
	}
	defer resp.Body.Close()

	var result WallUploadServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, "", 0, err
	}

	if result.Error.ErrorCode != 0 {
		return 0, "", 0, fmt.Errorf("VK API error: %s (code: %d)",
			result.Error.ErrorMsg, result.Error.ErrorCode)
	}

	return result.Response.AlbumID, result.Response.UploadURL, result.Response.UserID, nil
}

// uploadFileToServer uploads a file to the given server URL
func (c *Client) uploadFileToServer(uploadURL, filePath string) (int, string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, "", "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", filepath.Base(filePath))
	if err != nil {
		return 0, "", "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return 0, "", "", err
	}

	err = writer.Close()
	if err != nil {
		return 0, "", "", err
	}

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return 0, "", "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", "", fmt.Errorf("server returned non-200 status: %d", resp.StatusCode)
	}

	var result UploadPhotoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, "", "", err
	}

	return result.Server, result.Photo, result.Hash, nil
}

// saveWallPhoto saves the uploaded photo to the wall
func (c *Client) saveWallPhoto(photo string, server int, hash string) (int, int, error) {
	url := "https://api.vk.com/method/photos.saveWallPhoto"

	// Remove the minus sign from owner ID if present
	groupID := strings.Replace(c.OwnerID, "-", "", 1)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}

	q := req.URL.Query()
	q.Add("group_id", groupID)
	q.Add("access_token", c.token)
	q.Add("photo", photo)
	q.Add("server", strconv.Itoa(server))
	q.Add("hash", hash)
	q.Add("v", "5.199")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var result SaveWallPhotoResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return 0, 0, err
	}

	if result.Error.ErrorCode != 0 {
		return 0, 0, fmt.Errorf("VK API error: %s (code: %d)",
			result.Error.ErrorMsg, result.Error.ErrorCode)
	}

	if len(result.Response) == 0 {
		return 0, 0, fmt.Errorf("empty response from photos.saveWallPhoto")
	}

	return result.Response[0].OwnerID, result.Response[0].ID, nil
}

// videoSave gets the upload URL for videos
func (c *Client) videoSave() (int, string, error) {
	url := "https://api.vk.com/method/video.save"

	// Remove the minus sign from owner ID if present
	groupID := strings.Replace(c.OwnerID, "-", "", 1)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	q := req.URL.Query()
	q.Add("group_id", groupID)
	q.Add("access_token", c.token)
	q.Add("v", "5.199")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	var result VideoSaveResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return 0, "", err
	}

	if result.Error.ErrorCode != 0 {
		return 0, "", fmt.Errorf("VK API error: %s (code: %d)",
			result.Error.ErrorMsg, result.Error.ErrorCode)
	}

	return result.Response.OwnerID, result.Response.UploadURL, nil
}

// uploadVideoToServer uploads a video to the given server URL
func (c *Client) uploadVideoToServer(uploadURL, filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video_file", filepath.Base(filePath))
	if err != nil {
		return 0, 0, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return 0, 0, err
	}

	err = writer.Close()
	if err != nil {
		return 0, 0, err
	}

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	// The response format is somewhat inconsistent, so we need to handle different cases
	var result struct {
		OwnerID int    `json:"owner_id"`
		VideoID int    `json:"video_id"`
		Size    int    `json:"size"`
		Error   string `json:"error"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return 0, 0, err
	}

	if result.Error != "" {
		return 0, 0, fmt.Errorf("video upload error: %s", result.Error)
	}

	return result.OwnerID, result.VideoID, nil
}

// postToWall posts a message with attachments to the wall
func (c *Client) postToWall(message, attachments string, donutPaidDuration string) error {
	url := "https://api.vk.com/method/wall.post"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("access_token", c.token)
	q.Add("owner_id", c.OwnerID)
	q.Add("message", message)
	q.Add("attachments", attachments)
	q.Add("from_group", "1")
	q.Add("v", "5.199")
	if donutPaidDuration != "" { // if donutPaidDuration is presented
		q.Add("donut_paid_duration", donutPaidDuration)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result WallPostResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return err
	}

	if result.Error.ErrorCode != 0 {
		return fmt.Errorf("VK API error: %s (code: %d)",
			result.Error.ErrorMsg, result.Error.ErrorCode)
	}

	log.Printf("Posted successfully to wall, post ID: %d", result.Response.PostID)
	return nil
}
