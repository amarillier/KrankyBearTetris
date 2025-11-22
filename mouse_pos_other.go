//go:build !darwin
// +build !darwin

package main

import (
	"fyne.io/fyne/v2"
)

// getMousePosition returns the current mouse cursor position
// Stub implementation for non-macOS platforms
func getMousePosition() (x, y float32) {
	return 0, 0
}

// positionWindowOnMouseDisplay positions the window on the display containing the mouse cursor
// Fallback implementation for non-macOS platforms
func positionWindowOnMouseDisplay(app fyne.App, window fyne.Window) {
	// For non-macOS platforms, just center on screen
	window.CenterOnScreen()
}

