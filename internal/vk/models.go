package vk

// Client represents VK API client
type Client struct {
	AppID   string
	OwnerID string
	token   string
}

// WallUploadServerResponse represents the response from getWallUploadServer API
type WallUploadServerResponse struct {
	Response struct {
		AlbumID   int    `json:"album_id"`
		UploadURL string `json:"upload_url"`
		UserID    int    `json:"user_id"`
	} `json:"response"`
	Error struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error"`
}

// UploadPhotoResponse represents the response from photo upload API
type UploadPhotoResponse struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

// SaveWallPhotoResponse represents the response from photos.saveWallPhoto API
type SaveWallPhotoResponse struct {
	Response []struct {
		ID        int    `json:"id"` // Changed from string to int
		AlbumID   int    `json:"album_id"`
		OwnerID   int    `json:"owner_id"`
		UserID    int    `json:"user_id"`
		Sizes     []any  `json:"sizes"`
		Text      string `json:"text"`
		Date      int    `json:"date"`
		AccessKey string `json:"access_key"`
	} `json:"response"`
	Error struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error"`
}

// VideoSaveResponse represents the response from video.save API
type VideoSaveResponse struct {
	Response struct {
		UploadURL string `json:"upload_url"`
		VideoID   int    `json:"video_id"` // Changed from string to int
		OwnerID   int    `json:"owner_id"`
	} `json:"response"`
	Error struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error"`
}

// WallPostResponse represents the response from wall.post API
type WallPostResponse struct {
	Response struct {
		PostID int `json:"post_id"`
	} `json:"response"`
	Error struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error"`
}
