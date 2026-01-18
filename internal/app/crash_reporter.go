package app

import (
	"HyLauncher/internal/diagnostics"
	"HyLauncher/internal/env"
	"encoding/json"
	"os"
	"path/filepath"
)

// GetCrashReports returns all crash reports
func (a *App) GetCrashReports() ([]diagnostics.CrashReport, error) {
	crashDir := filepath.Join(env.GetDefaultAppDir(), "crashes")

	entries, err := os.ReadDir(crashDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []diagnostics.CrashReport{}, nil
		}
		return nil, err
	}

	var reports []diagnostics.CrashReport
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := os.ReadFile(filepath.Join(crashDir, entry.Name()))
		if err != nil {
			continue
		}

		var report diagnostics.CrashReport
		if err := json.Unmarshal(data, &report); err != nil {
			continue
		}

		reports = append(reports, report)
	}

	return reports, nil
}
