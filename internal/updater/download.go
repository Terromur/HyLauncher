package updater

import (
	"HyLauncher/pkg/download"
	"HyLauncher/pkg/fileutil"
	"context"
	"fmt"
	"os"
)

// Downloads latest launcher, returns path to temp file. If cant download deletes temp file
func DownloadTemp(
	ctx context.Context,
	url string,
	progress func(stage string, progress float64, message string, currentFile string, speed string, downloaded, total int64),
) (string, error) {

	tmpPath, err := fileutil.CreateTempFile("file-update-*")
	if err != nil {
		return "", err
	}

	if err := download.DownloadWithProgress(tmpPath, url, "update", 1.0, progress); err != nil {
		_ = os.Remove(tmpPath)
		return "", err
	}

	fmt.Printf("Download complete: %s\n", tmpPath)
	return tmpPath, nil
}
