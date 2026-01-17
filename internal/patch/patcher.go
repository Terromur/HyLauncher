package patch

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"HyLauncher/internal/env"
	"HyLauncher/internal/platform"
	"HyLauncher/pkg/download"
)

func ApplyPWR(ctx context.Context, pwrFile string,
	progressCallback func(stage string, progress float64, message string, currentFile string, speed string, downloaded, total int64)) error {

	gameLatest := filepath.Join(env.GetDefaultAppDir(), "release", "package", "game", "latest")
	stagingDir := filepath.Join(gameLatest, "staging-temp")
	_ = os.MkdirAll(stagingDir, 0755)

	butlerPath := filepath.Join(env.GetDefaultAppDir(), "tools", "butler", "butler")
	if runtime.GOOS == "windows" {
		butlerPath += ".exe"
	}

	cmd := exec.CommandContext(ctx, butlerPath,
		"apply",
		"--staging-dir", stagingDir,
		pwrFile,
		gameLatest,
	)

	platform.HideConsoleWindow(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if progressCallback != nil {
		progressCallback("game", 60, "Applying game patch...", "", "", 0, 0)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	// Retry rename on Windows if locked
	if runtime.GOOS == "windows" {
		for i := 0; i < 5; i++ {
			if err := os.Rename(stagingDir, gameLatest); err == nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
	} else {
		_ = os.Rename(stagingDir, gameLatest)
	}

	if progressCallback != nil {
		progressCallback("game", 100, "Game installed successfully", "", "", 0, 0)
	}

	return nil
}

func DownloadPWR(ctx context.Context, versionType string, prevVer int, targetVer int,
	progressCallback func(stage string, progress float64, message string, currentFile string, speed string, downloaded, total int64)) (string, error) {

	cacheDir := filepath.Join(env.GetDefaultAppDir(), "cache")
	_ = os.MkdirAll(cacheDir, 0755)

	osName := runtime.GOOS
	arch := runtime.GOARCH

	fileName := fmt.Sprintf("%d.pwr", targetVer)
	dest := filepath.Join(cacheDir, fileName)
	tempDest := dest + ".tmp"

	_ = os.Remove(tempDest)

	if _, err := os.Stat(dest); err == nil {
		if progressCallback != nil {
			progressCallback("game", 40, "PWR file cached", fileName, "", 0, 0)
		}
		return dest, nil
	}

	url := fmt.Sprintf("https://game-patches.hytale.com/patches/%s/%s/%s/%d/%s",
		osName, arch, versionType, prevVer, fileName)

	if err := download.DownloadWithProgress(tempDest, url, "game", 0.4, progressCallback); err != nil {
		_ = os.Remove(tempDest)
		return "", err
	}

	if err := os.Rename(tempDest, dest); err != nil {
		_ = os.Remove(tempDest)
		return "", err
	}

	return dest, nil
}
