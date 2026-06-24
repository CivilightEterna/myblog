package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"blog-admin-api/database"
	"blog-admin-api/models"
)

type BuildService struct {
	BlogRoot string
}

func NewBuildService(blogRoot string) *BuildService {
	return &BuildService{BlogRoot: blogRoot}
}

// toBashPath converts a Windows path to MSYS/Cygwin bash-compatible path
func toBashPath(p string) string {
	// E:\code\blog -> /e/code/blog
	p = strings.ReplaceAll(p, "\\", "/")
	if len(p) >= 2 && p[1] == ':' {
		drive := strings.ToLower(string(p[0]))
		p = "/" + drive + p[2:]
	}
	return p
}

func (s *BuildService) ExecuteBuild(triggerType string) (*models.BuildRecord, error) {
	record := models.BuildRecord{
		Status:      "running",
		TriggerType: triggerType,
		StartedAt:   time.Now(),
	}
	database.DB.Create(&record)

	// Build script path — convert to bash-compatible path
	scriptPath := filepath.Join(s.BlogRoot, "scripts", "build.sh")
	bashScriptPath := toBashPath(scriptPath)

	var logBuf bytes.Buffer

	// Run build script
	cmd := exec.Command("bash", bashScriptPath)
	cmd.Dir = toBashPath(s.BlogRoot)
	cmd.Stdout = &logBuf
	cmd.Stderr = &logBuf
	cmd.Env = append(os.Environ(), "BLOG_ROOT="+toBashPath(s.BlogRoot))

	err := cmd.Run()
	now := time.Now()
	record.FinishedAt = &now

	if err != nil {
		record.Status = "failed"
		record.Log = fmt.Sprintf("Build failed: %s\n\n%s", err.Error(), logBuf.String())
	} else {
		record.Status = "success"
		record.Log = logBuf.String()

		// Determine release path
		record.ReleasePath = filepath.Join(s.BlogRoot, "releases", now.Format("20060102_150405"))
	}

	database.DB.Save(&record)
	return &record, err
}
