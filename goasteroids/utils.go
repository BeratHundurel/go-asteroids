package goasteroids

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// sharedDrawOp is a reusable DrawImageOptions that eliminates per-frame heap
// allocations in entity Draw methods. Always call sharedDrawOp.GeoM.Reset()
// before use. Safe because ebiten copies all options at DrawImage call time
// and all Draw calls happen on a single goroutine.
var sharedDrawOp = &ebiten.DrawImageOptions{}

// sharedColormDrawOp is the colorm equivalent of sharedDrawOp.
var sharedColormDrawOp = &colorm.DrawImageOptions{}

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

// separateMeteors resolves overlaps between a slice of meteors by pushing
// each overlapping pair apart and reflecting their velocities along the
// collision normal (elastic collision). The physics body position is also
// synced after every push, keeping resolv in step with the logical position.
func separateMeteors(meteors []*Meteor) {
	n := len(meteors)
	if n < 2 {
		return
	}

	for i := range n {
		m1 := meteors[i]
		r1 := m1.radius
		for j := i + 1; j < n; j++ {
			m2 := meteors[j]

			dx := m2.position.X - m1.position.X
			dy := m2.position.Y - m1.position.Y
			distSq := dx*dx + dy*dy
			minDist := r1 + m2.radius

			if distSq == 0 || distSq >= minDist*minDist {
				continue
			}

			dist := math.Sqrt(distSq)
			// Collision normal
			nx, ny := dx/dist, dy/dist

			// Push both meteors apart so they no longer overlap
			overlap := (minDist - dist) / 2
			m1.position.X -= nx * overlap
			m1.position.Y -= ny * overlap
			m2.position.X += nx * overlap
			m2.position.Y += ny * overlap
			m1.meteorObj.SetPosition(m1.position.X, m1.position.Y)
			m2.meteorObj.SetPosition(m2.position.X, m2.position.Y)

			// Elastic collision: exchange velocity components along the collision normal
			dvx := m1.movement.X - m2.movement.X
			dvy := m1.movement.Y - m2.movement.Y
			dot := dvx*nx + dvy*ny

			// Only swap if the meteors are actually moving toward each other
			if dot > 0 {
				m1.movement.X -= dot * nx
				m1.movement.Y -= dot * ny
				m2.movement.X += dot * nx
				m2.movement.Y += dot * ny
			}
		}
	}
}
