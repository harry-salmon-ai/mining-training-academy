package handlers

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func mediaDir() string {
	if _, err := os.Stat("/mnt/data"); err == nil {
		dir := "/mnt/data/images"
		os.MkdirAll(dir, 0755)
		return dir
	}
	// Local fallback
	dir := "./uploads/images"
	os.MkdirAll(dir, 0755)
	return dir
}

var allowedMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
	"image/svg+xml": ".svg",
}

const maxUploadSize = 10 << 20 // 10 MB

func UploadMedia(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)
	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 10 MB)"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Detect MIME from first 512 bytes
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	detectedMIME := http.DetectContentType(buf[:n])
	file.Seek(0, io.SeekStart)

	// Also check Content-Type header as fallback for SVG (DetectContentType can't identify SVG)
	headerMIME := header.Header.Get("Content-Type")
	ext, ok := allowedMIME[detectedMIME]
	if !ok {
		ext, ok = allowedMIME[headerMIME]
	}
	// Try by file extension
	if !ok {
		extFromName := strings.ToLower(filepath.Ext(header.Filename))
		for mimeType, mimeExt := range allowedMIME {
			if mimeExt == extFromName {
				ext = mimeExt
				_ = mimeType
				ok = true
				break
			}
		}
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported file type: %s", detectedMIME)})
		return
	}

	filename := uuid.New().String() + ext
	destPath := filepath.Join(mediaDir(), filename)

	dst, err := os.Create(destPath)
	if err != nil {
		log.Printf("Failed to create media file %s: %v", destPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	url := "/api/media/" + filename
	c.JSON(http.StatusCreated, gin.H{
		"url":      url,
		"filename": filename,
	})
}

func ServeMedia(c *gin.Context) {
	filename := c.Param("filename")
	// Prevent path traversal
	filename = filepath.Base(filename)

	filePath := filepath.Join(mediaDir(), filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.File(filePath)
}
