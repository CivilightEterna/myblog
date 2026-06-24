package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"blog-admin-api/database"
	"blog-admin-api/models"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	BlogRoot string
}

func NewUploadHandler(blogRoot string) *UploadHandler {
	return &UploadHandler{BlogRoot: blogRoot}
}

// magicBytes maps MIME types to their expected file header magic bytes
var allowedMagicBytes = map[string][]byte{
	"image/jpeg": {0xFF, 0xD8, 0xFF},
	"image/png":  {0x89, 0x50, 0x4E, 0x47},
	"image/gif":  {0x47, 0x49, 0x46, 0x38},
	"image/webp": {0x52, 0x49, 0x46, 0x46},
}

var allowedExts = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}
	defer file.Close()

	// === 1. Path traversal prevention ===
	originalName := header.Filename
	if strings.Contains(originalName, "..") ||
		strings.Contains(originalName, "/") ||
		strings.Contains(originalName, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名包含非法字符"})
		return
	}

	// === 2. Extension whitelist ===
	ext := strings.ToLower(filepath.Ext(originalName))
	expectedMime, ok := allowedExts[ext]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型，允许: jpg, png, gif, webp, svg"})
		return
	}

	// === 3. File size limit (5MB) ===
	const maxSize = 5 * 1024 * 1024
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 5MB"})
		return
	}

	// === 4. Magic byte validation (prevent extension spoofing) ===
	if ext != ".svg" {
		// Read first 512 bytes for magic byte check
		headerBytes := make([]byte, 512)
		n, _ := io.ReadFull(file, headerBytes)
		if n == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取文件内容"})
			return
		}

		matched := false
		for mime, magic := range allowedMagicBytes {
			if len(headerBytes) >= len(magic) && compareBytes(headerBytes[:len(magic)], magic) {
				matched = true
				// Override MIME type from magic bytes (more reliable than Content-Type header)
				expectedMime = mime
				break
			}
		}
		if !matched {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件内容与扩展名不匹配"})
			return
		}

		// Seek back to beginning for saving
		file.Seek(0, 0)
	}

	// === 5. SVG security: log a warning ===
	if ext == ".svg" {
		fmt.Printf("⚠ SVG file uploaded: %s — ensure SVG does not contain scripts\n", originalName)
	}

	// === 6. Generate safe filename ===
	now := time.Now()
	yearMonth := now.Format("2006/01")
	file.Seek(0, 0)
	hash := sha256.New()
	io.Copy(hash, file)
	fileHash := fmt.Sprintf("%x", hash.Sum(nil))[:16]
	newFilename := fmt.Sprintf("%s_%s%s", now.Format("20060102_150405"), fileHash, ext)

	// === 7. Save to uploads directory ===
	uploadsDir := filepath.Join(h.BlogRoot, "apps", "blog-web", "public", "uploads", yearMonth)
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	// Resolve the final path and verify it's inside the uploads directory
	dstPath := filepath.Join(uploadsDir, newFilename)
	absDst, _ := filepath.Abs(dstPath)
	absUploads, _ := filepath.Abs(uploadsDir)
	if !strings.HasPrefix(absDst, absUploads) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法文件路径"})
		return
	}

	file.Seek(0, 0)
	dst, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入文件失败"})
		return
	}

	// === 8. Record in database ===
	url := fmt.Sprintf("/uploads/%s/%s", yearMonth, newFilename)
	upload := models.UploadFile{
		Filename:     newFilename,
		OriginalName: originalName,
		Path:         dstPath,
		URL:          url,
		MimeType:     expectedMime,
		Size:         header.Size,
		Hash:         fileHash,
	}

	if err := database.DB.Create(&upload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录上传文件失败"})
		return
	}

	c.JSON(http.StatusCreated, upload)
}

func (h *UploadHandler) ListUploads(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.UploadFile{}).Count(&total)

	var uploads []models.UploadFile
	if err := database.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&uploads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     uploads,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *UploadHandler) DeleteUpload(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件ID"})
		return
	}

	var upload models.UploadFile
	if err := database.DB.First(&upload, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	if err := os.Remove(upload.Path); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
		return
	}

	database.DB.Delete(&upload)
	c.JSON(http.StatusOK, gin.H{"message": "文件已删除"})
}

func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
