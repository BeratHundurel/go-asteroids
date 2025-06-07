package goasteroids

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func HalfOfTheImage(image *ebiten.Image) (float64, float64) {
	bounds := image.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	return halfW, halfH
}

func getAppDataDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(u.HomeDir, "Library", "Application Support", "Go Asteroids"), nil
	case "windows":
		return filepath.Join(u.HomeDir, "AppData", "Roaming", "Go Asteroids"), nil
	default:
		return filepath.Join(u.HomeDir, ".local", "share", "Go Asteroids"), nil
	}
}

func getHighScore() (int, error) {
	dir, err := getAppDataDir()
	if err != nil {
		return 0, fmt.Errorf("failed to get app data directory: %w", err)
	}

	// Create directory with all parent directories
	if err := os.MkdirAll(dir, 0750); err != nil {
		return 0, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	scoreFile := filepath.Join(dir, "highscore.txt")
	if _, err := os.Stat(scoreFile); os.IsNotExist(err) {
		err := os.WriteFile(scoreFile, []byte("0"), 0644)
		if err != nil {
			return 0, fmt.Errorf("failed to create highscore file: %w", err)
		}
	}

	contents, err := os.ReadFile(scoreFile)
	if err != nil {
		return 0, fmt.Errorf("failed to read highscore file: %w", err)
	}

	score := strings.TrimSpace(string(contents))
	s, err := strconv.Atoi(score)
	if err != nil {
		return 0, fmt.Errorf("failed to convert highscore to integer: %w", err)
	}
	return s, nil
}

func updateHighScore(score int) error {
	dir, err := getAppDataDir()
	if err != nil {
		return fmt.Errorf("failed to get app data directory: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	scoreFile := filepath.Join(dir, "highscore.txt")
	s := fmt.Sprintf("%d", score)
	err = os.WriteFile(scoreFile, []byte(s), 0644)
	if err != nil {
		return fmt.Errorf("failed to update highscore file: %w", err)
	}
	return nil
}
