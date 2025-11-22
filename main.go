package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

const (
	// appName    = "KrankyBear Tetris"
	appVersion = "0.1.1" // see FyneApp.toml
	appAuthor  = "Allan Marillier"
)

var appName = "KrankyBear Tetris"
var appCopyright = "Copyright (c) Allan Marillier, 2025-" + strconv.Itoa(time.Now().Year())

func main() {
	// Create Fyne app
	myApp := app.NewWithID("com.github.amarillier.KrankyBearTetris")
	
	// Set app icon (for taskbar on Windows)
	myApp.SetIcon(resourceKrankyBearVikingHelmetPng)
	
	// Create game instance
	game := NewGame()
	
	// Create UI
	ui := NewGameUI(myApp, game)
	
	// Set up game loop
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
		defer ticker.Stop()
		
		for range ticker.C {
			game.Update()
			// Wrap UI updates in fyne.Do() for thread-safety when called from goroutine
			fyne.Do(func() {
				canvas.Refresh(ui.boardCanvas)
				canvas.Refresh(ui.nextCanvas)
				ui.updateLabelsOnly()
			})
		}
	}()
	
	// Show window and run app
	ui.ShowAndRun()
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
