package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var ConfPath string
var StartTime = time.Now()

// Get the currently selected configuration file directory
// For non-Windows systems, select the /etc/nps as config directory if exist, or select ./
// windows system, select the C:\Program Files\nps as config directory if exist, or select ./
func GetRunPath() string {
	var path string
	if len(os.Args) == 1 {
		if !IsWindows() {
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) //返回
			return dir + "/"
		} else {
			return "./"
		}
	} else {
		if path = GetInstallPath(); !FileExists(path) {
			return GetAppPath()
		}
	}
	return path
}

// Different systems get different installation paths
func GetInstallPath() string {
	var path string

	if ConfPath != "" {
		return ConfPath
	}

	if IsWindows() {
		path = `C:\Program Files\nps`
	} else {
		path = "/etc/nps"
	}

	return path
}

// Get the absolute path to the running directory
func GetAppPath() string {
	if exePath, err := os.Executable(); err == nil {
		return filepath.Dir(exePath)
	}
	if path, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return path
	}
	return os.Args[0]
}

// Determine whether the current system is a Windows system?
func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

// interface log file path
func GetLogPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), "nps.log")
	} else {
		path = "/var/log/nps.log"
	}
	return path
}

// interface npc log file path
func GetNpcLogPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), "npc.log")
	} else {
		path = "/var/log/npc.log"
	}
	return path
}

// interface pid file path
func GetTmpPath() string {
	var path string
	if IsWindows() {
		path = GetAppPath()
	} else {
		path = "/tmp"
	}
	return path
}

// config file path
func GetConfigPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), "conf/npc.conf")
	} else {
		path = "conf/npc.conf"
	}
	return path
}

func ResolvePath(path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(GetRunPath(), path)
	}
	return path
}

func GetRunTime() string {
	totalSecs := int64(time.Since(StartTime).Seconds())
	days := totalSecs / 86400
	totalSecs %= 86400
	hours := totalSecs / 3600
	totalSecs %= 3600
	mins := totalSecs / 60
	secs := totalSecs % 60
	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if mins > 0 {
		parts = append(parts, fmt.Sprintf("%dm", mins))
	}
	if secs > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", secs))
	}
	return strings.Join(parts, " ")
}

func GetRunSecs() int64 {
	return int64(time.Since(StartTime).Seconds())
}

func GetStartTime() int64 {
	return StartTime.Unix()
}
