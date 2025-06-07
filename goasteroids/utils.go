package goasteroids

import (
	"fmt"
	"os"
	"os/user"
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

func getHighScore() (int, error) {
	u, err := user.Current()
	if err != nil {
		return 0, err
	}
	path := ""
	switch runtime.GOOS {
	case "darwin":
		path = fmt.Sprintf("/Users/%s/Library/Application Support/Go Asteroids",
			u.Username)
	case "windows":
		path = fmt.Sprintf("C:\\Users\\%s\\AppData\\Go Asteroids",
			u.Username)
	default:
		path = fmt.Sprintf("/users/%s", u.Username)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0750); err != nil {
			return 0, fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}

	if _, err := os.Stat(path + "/highscore.txt"); os.IsNotExist(err) {
		err := os.WriteFile(path+"/highscore.txt", []byte("0"), 0750)
		if err != nil {
			return 0, fmt.Errorf("failed to create highscore file: %w", err)
		}
	}

	contents, err := os.ReadFile(path + "/highscore.txt")
	if err != nil {
		return 0, fmt.Errorf("failed to read highscore file: %w", err)
	}

	score := string(contents)
	score = strings.TrimSpace(score)
	s, err := strconv.Atoi(score)
	if err != nil {
		return 0, fmt.Errorf("failed to convert highscore to integer: %w", err)
	}
	return s, nil
}

func updateHighScore(score int) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	path := ""
	switch runtime.GOOS {
	case "darwin":
		path = fmt.Sprintf("/Users/%s/Library/Application Support/Go Asteroids/highscore.txt",
			u.Username)
	case "windows":
		path = fmt.Sprintf("C:\\Users\\%s\\AppData\\Roaming\\Go Asteroids\\highscore.txt",
			u.Username)
	default:
		path = fmt.Sprintf("/users/%s/highscore.txt", u.Username)
	}

	s := fmt.Sprintf("%d", score)
	err = os.WriteFile(path+"/highscore.txt", []byte(s), 0750)
	if err != nil {
		return fmt.Errorf("failed to update highscore file: %w", err)
	}
	return nil
}
