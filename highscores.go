package main

import (
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
)

const maxHighScores = 5

// HighScoreManager manages high scores using Fyne preferences
type HighScoreManager struct {
	prefs fyne.Preferences
}

// NewHighScoreManager creates a new high score manager
func NewHighScoreManager(app fyne.App) *HighScoreManager {
	return &HighScoreManager{
		prefs: app.Preferences(),
	}
}

// GetHighScores retrieves the high scores from preferences
func (h *HighScoreManager) GetHighScores() []int {
	scores := make([]int, 0, maxHighScores)

	for i := 0; i < maxHighScores; i++ {
		key := "highscore_" + strconv.Itoa(i)
		score := h.prefs.IntWithFallback(key, 0)
		if score > 0 {
			scores = append(scores, score)
		}
	}

	// Sort in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	return scores
}

// AddScore adds a score if it's high enough, maintaining top 5
// Returns true if this score is a new high score (in top 5)
func (h *HighScoreManager) AddScore(score int) bool {
	if score <= 0 {
		return false
	}

	scores := h.GetHighScores()

	// Check if score qualifies for top 5
	isHighScore := len(scores) < maxHighScores || (len(scores) > 0 && score > scores[len(scores)-1])

	if isHighScore {
		scores = append(scores, score)

		// Sort in descending order
		sort.Sort(sort.Reverse(sort.IntSlice(scores)))

		// Keep only top maxHighScores
		if len(scores) > maxHighScores {
			scores = scores[:maxHighScores]
		}

		// Save back to preferences
		for i := 0; i < maxHighScores; i++ {
			key := "highscore_" + strconv.Itoa(i)
			if i < len(scores) {
				h.prefs.SetInt(key, scores[i])
			} else {
				h.prefs.SetInt(key, 0)
			}
		}

		return true
	}

	return false
}

// ResetHighScores clears all high scores
func (h *HighScoreManager) ResetHighScores() {
	for i := 0; i < maxHighScores; i++ {
		key := "highscore_" + strconv.Itoa(i)
		h.prefs.SetInt(key, 0)
	}
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
