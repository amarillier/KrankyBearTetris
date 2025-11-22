package main

import (
	"image"
	"image/color"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	updatechecker "github.com/amarillier/go-update-checker"
)

// ShortcutBossKey describes the boss key shortcut (F12) - pause and hide window
// Using F12 as it's commonly used for boss keys and less likely to conflict
type ShortcutBossKey struct{}

var _ fyne.KeyboardShortcut = (*ShortcutBossKey)(nil)

// Key returns the KeyName for this shortcut
func (s *ShortcutBossKey) Key() fyne.KeyName {
	return fyne.KeyF12
}

// Mod returns the KeyModifier for this shortcut
func (s *ShortcutBossKey) Mod() fyne.KeyModifier {
	return fyne.KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (s *ShortcutBossKey) ShortcutName() string {
	return "BossKey"
}

const (
	BlockSize = 25
	BoardX    = 50
	BoardY    = 50
)

// GameUI manages the game's user interface
type GameUI struct {
	game             *Game
	app              fyne.App
	window           fyne.Window
	boardCanvas      *canvas.Raster
	nextCanvas       *canvas.Raster
	scoreLabel       *widget.Label
	linesLabel       *widget.Label
	levelLabel       *widget.Label
	statusLabel      *widget.Label
	highScoreLabel   *widget.Label
	startButton      *widget.Button
	pauseButton      *widget.Button
	resetScoreButton *widget.Button
	soundButton      *widget.Button
	highScoreMgr     *HighScoreManager
	soundMgr         *SoundManager
	resetDialog      fyne.Window // Track reset confirmation dialog window
	aboutDialog      fyne.Window // Track about dialog window
	helpDialog       fyne.Window // Track help dialog window
	updateDialog     fyne.Window // Track update check dialog window
	showMenuItem     *fyne.MenuItem
	hideMenuItem     *fyne.MenuItem
	windowVisible    bool   // Track window visibility state
	versionStatus    string // Track version status: "current", "newer", or "unknown"
}

// NewGameUI creates a new game UI
func NewGameUI(app fyne.App, game *Game) *GameUI {
	ui := &GameUI{
		game:          game,
		app:           app,
		window:        app.NewWindow(appName),
		highScoreMgr:  NewHighScoreManager(app),
		soundMgr:      NewSoundManager(),
		windowVisible: true, // Window starts visible
	}

	ui.setupUI()
	ui.setupKeyboard()
	ui.setupMenu()

	// Set up sound callback
	game.SetSoundCallback(ui.handleSoundEvent)

	// Check for updates on startup
	ui.checkForUpdates()

	// Setup system tray after everything else is ready
	ui.setupSystemTray()

	return ui
}

// updateChecker checks for version updates
func (ui *GameUI) updateChecker() (string, bool) {
	uc := updatechecker.New("amarillier", "KrankyBearTetris", "KrankyBear Tetris", "https://github.com/amarillier/KrankyBearTetris/releases/latest", 0, false)
	uc.CheckForUpdate(appVersion)
	return uc.Message, uc.UpdateAvailable
}

// checkForUpdates checks for updates and updates version status
func (ui *GameUI) checkForUpdates() {
	updtmsg, updateAvailable := ui.updateChecker()
	// Check if we're running a newer version than released
	if strings.Contains(updtmsg, "running a newer version") || strings.Contains(updtmsg, "newer than") {
		ui.versionStatus = "newer" // We're running newer than released
	} else if strings.Contains(updtmsg, "running the latest") || strings.Contains(updtmsg, "up to date") {
		ui.versionStatus = "current" // Version matches released
	} else if updateAvailable {
		ui.versionStatus = "update" // Update available
	} else {
		ui.versionStatus = "unknown"
	}
}

// setupUI creates the UI layout
func (ui *GameUI) setupUI() {
	ui.window.Resize(fyne.NewSize(500, 700))
	// Center window on screen - Fyne will center on the primary display
	ui.window.CenterOnScreen()

	// Create canvas for game board
	ui.boardCanvas = canvas.NewRaster(ui.drawBoard)
	boardSize := fyne.NewSize(BoardWidth*BlockSize, BoardHeight*BlockSize)
	ui.boardCanvas.Resize(boardSize)
	ui.boardCanvas.SetMinSize(boardSize)

	// Create canvas for next piece preview
	ui.nextCanvas = canvas.NewRaster(ui.drawNextPiece)
	nextSize := fyne.NewSize(4*BlockSize, 4*BlockSize)
	ui.nextCanvas.Resize(nextSize)
	ui.nextCanvas.SetMinSize(nextSize)

	// Create labels
	ui.scoreLabel = widget.NewLabel("Score: 0")
	ui.linesLabel = widget.NewLabel("Lines: 0")
	ui.levelLabel = widget.NewLabel("Level: 1")
	ui.statusLabel = widget.NewLabel("Press Start to begin")
	ui.statusLabel.Alignment = fyne.TextAlignCenter
	ui.highScoreLabel = widget.NewLabel("High Scores:\n" + ui.formatHighScores())

	// Create buttons
	ui.startButton = widget.NewButton("Start", ui.onStart)
	ui.pauseButton = widget.NewButton("Pause", ui.onPause)
	ui.pauseButton.Disable()
	ui.resetScoreButton = widget.NewButton("Reset Scores", ui.onResetScores)
	ui.soundButton = widget.NewButton("🔊 Sound On", ui.onToggleSound)
	ui.updateSoundButton()

	// Layout
	infoPanel := container.NewVBox(
		widget.NewLabel("Next:"),
		ui.nextCanvas,
		widget.NewSeparator(),
		ui.scoreLabel,
		ui.linesLabel,
		ui.levelLabel,
		widget.NewSeparator(),
		ui.statusLabel,
		widget.NewSeparator(),
		ui.highScoreLabel,
		widget.NewSeparator(),
		ui.startButton,
		ui.pauseButton,
		ui.resetScoreButton,
		ui.soundButton,
		widget.NewLabel("\nControls:\n← → Move\n↓ Soft Drop\n↑ Rotate\nSpace Hard Drop\nF12 Boss Key (Pause & Hide)"),
	)

	// Use HBox to prevent centering and extra space
	content := container.NewBorder(
		nil, nil, nil, infoPanel,
		ui.boardCanvas, // Direct canvas, no centering container
	)

	ui.window.SetContent(content)
	// Force initial refresh of canvas
	ui.boardCanvas.Refresh()
	ui.nextCanvas.Refresh()

	// Intercept window close (X button) to quit the application
	ui.window.SetCloseIntercept(func() {
		ui.app.Quit()
	})

	ui.window.SetOnClosed(func() {
		// Cleanup if needed
	})
}

// setupMenu sets up the menu bar
func (ui *GameUI) setupMenu() {

	showMenu := fyne.NewMenuItem("Show", func() {
		fyne.Do(ui.showWindow)
	})
	hideMenu := fyne.NewMenuItem("Hide", func() {
		fyne.Do(ui.hideWindow)
	})
	aboutMenu := fyne.NewMenuItem("About", func() {
		fyne.Do(ui.showAboutDialog)
	})
	helpMenu := fyne.NewMenuItem("Help", func() {
		fyne.Do(ui.showHelpDialog)
	})
	updateMenu := fyne.NewMenuItem("Check for Update", func() {
		ui.showUpdateDialog()
	})

	mainMenu := fyne.NewMenu("KrankyBear Tetris",
		showMenu,
		hideMenu,
		fyne.NewMenuItemSeparator(),
		aboutMenu,
		helpMenu,
		updateMenu,
	)

	menu := fyne.NewMainMenu(mainMenu)
	ui.window.SetMainMenu(menu)
}

// setupSystemTray sets up the system tray icon and menu using desktop.App interface
func (ui *GameUI) setupSystemTray() {
	// Check if app supports desktop features (system tray)
	if desk, ok := ui.app.(desktop.App); ok {
		// Create menu items
		ui.showMenuItem = fyne.NewMenuItem("Show", func() {
			fyne.Do(ui.showWindow)
		})
		ui.hideMenuItem = fyne.NewMenuItem("Hide", func() {
			fyne.Do(ui.hideWindow)
		})
		fyne.NewMenuItemSeparator()
		helpTray := fyne.NewMenuItem("Help", func() {
			fyne.Do(ui.showHelpDialog)
		})
		aboutTray := fyne.NewMenuItem("About", func() {
			fyne.Do(ui.showAboutDialog)
		})
		updateTray := fyne.NewMenuItem("Check for Update", func() {
			ui.showUpdateDialog()
		})
		fyne.NewMenuItemSeparator()
		quitTray := fyne.NewMenuItem("Quit", func() {
			fyne.Do(ui.app.Quit)
		})

		// Create menu
		menu := fyne.NewMenu("KrankyBear Tetris",
			ui.showMenuItem,
			ui.hideMenuItem,
			fyne.NewMenuItemSeparator(),
			aboutTray,
			helpTray,
			updateTray,
			fyne.NewMenuItemSeparator(),
			quitTray,
		)

		// Set system tray menu and icon
		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(resourceKrankyBearVikingHelmetPng)

		// Update menu state based on window visibility
		ui.updateTrayMenuState()
	}
}

// updateTrayMenuState updates the tray menu items based on window visibility
func (ui *GameUI) updateTrayMenuState() {
	// Recreate the menu with updated state
	if desk, ok := ui.app.(desktop.App); ok {
		// Create menu items
		aboutTray := fyne.NewMenuItem("About", func() {
			fyne.Do(ui.showAboutDialog)
		})
		helpTray := fyne.NewMenuItem("Help", func() {
			fyne.Do(ui.showHelpDialog)
		})
		updateTray := fyne.NewMenuItem("Check for Update", func() {
			ui.showUpdateDialog()
		})

		// Create show/hide menu items - always create both, but we'll handle logic in callbacks
		ui.showMenuItem = fyne.NewMenuItem("Show", func() {
			fyne.Do(func() {
				if !ui.windowVisible {
					ui.showWindow()
				}
			})
		})
		ui.hideMenuItem = fyne.NewMenuItem("Hide", func() {
			fyne.Do(func() {
				if ui.windowVisible {
					ui.hideWindow()
				}
			})
		})

		quitTray := fyne.NewMenuItem("Quit", func() {
			fyne.Do(ui.app.Quit)
		})

		// Create menu
		menu := fyne.NewMenu("KrankyBear Tetris",
			aboutTray,
			helpTray,
			updateTray,
			fyne.NewMenuItemSeparator(),
			ui.showMenuItem,
			ui.hideMenuItem,
			fyne.NewMenuItemSeparator(),
			quitTray,
		)

		// Update system tray menu
		desk.SetSystemTrayMenu(menu)
	}
}

// showWindow shows the main window
func (ui *GameUI) showWindow() {
	ui.window.Show()
	ui.window.RequestFocus()
	ui.windowVisible = true
	ui.updateTrayMenuState()
}

// hideWindow hides the main window
func (ui *GameUI) hideWindow() {
	ui.window.Hide()
	ui.windowVisible = false
	ui.updateTrayMenuState()
}

// setupKeyboard sets up keyboard controls
func (ui *GameUI) setupKeyboard() {
	ui.window.Canvas().AddShortcut(&fyne.ShortcutCopy{}, func(shortcut fyne.Shortcut) {
		// Handle shortcuts if needed
	})

	// Add boss key shortcut (F12) - pause game and hide window
	// Using F12 as it's commonly used for boss keys and less likely to conflict
	ui.window.Canvas().AddShortcut(&ShortcutBossKey{}, func(shortcut fyne.Shortcut) {
		ui.bossKey()
	})

	// Also handle F12 in key handler as backup
	ui.window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		// Check for boss key (F12)
		if ev.Name == fyne.KeyF12 {
			ui.bossKey()
			return
		}

		// Handle other keys normally
		ui.HandleKey(ev.Name)
	})
}

// drawBoard draws the game board
func (ui *GameUI) drawBoard(width, height int) image.Image {
	// Calculate exact board dimensions
	boardPixelWidth := BoardWidth * BlockSize
	boardPixelHeight := BoardHeight * BlockSize

	// Use board dimensions, not the canvas size (which might be larger)
	imgWidth := boardPixelWidth
	imgHeight := boardPixelHeight
	if width > 0 && width < imgWidth {
		imgWidth = width
	}
	if height > 0 && height < imgHeight {
		imgHeight = height
	}

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Draw background only within board area
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			img.Set(x, y, color.RGBA{20, 20, 30, 255})
		}
	}

	// Draw grid - only within board boundaries
	gridColor := color.RGBA{40, 40, 50, 255}

	// Draw horizontal grid lines (between rows)
	for y := 0; y <= BoardHeight; y++ {
		py := y * BlockSize
		if py < imgHeight {
			for x := 0; x < imgWidth; x++ {
				img.Set(x, py, gridColor)
			}
		}
	}

	// Draw vertical grid lines (between columns)
	for x := 0; x <= BoardWidth; x++ {
		px := x * BlockSize
		if px < imgWidth {
			for y := 0; y < imgHeight; y++ {
				img.Set(px, y, gridColor)
			}
		}
	}

	// Draw placed blocks
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if ui.game.Board[y][x] != nil {
				c := ui.game.Board[y][x]
				ui.drawBlock(img, x*BlockSize, y*BlockSize, BlockSize, color.RGBA{c.R, c.G, c.B, 255})
			}
		}
	}

	// Draw current piece
	if ui.game.CurrentPiece != nil && (ui.game.State == StatePlaying || ui.game.State == StatePaused) {
		shape := ui.game.CurrentPiece.GetRotatedShape()
		pieceColor := TetrominoColors[ui.game.CurrentPiece.Type]
		c := color.RGBA{pieceColor.R, pieceColor.G, pieceColor.B, 255}

		for y := 0; y < len(shape); y++ {
			for x := 0; x < len(shape[y]); x++ {
				if shape[y][x] {
					boardX := ui.game.CurrentPiece.X + x
					boardY := ui.game.CurrentPiece.Y + y

					if boardY >= 0 && boardX >= 0 && boardX < BoardWidth {
						px := boardX * BlockSize
						py := boardY * BlockSize
						if py >= 0 && py < height {
							ui.drawBlock(img, px, py, BlockSize, c)
						}
					}
				}
			}
		}
	}

	return img
}

// drawNextPiece draws the next piece preview
func (ui *GameUI) drawNextPiece(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{20, 20, 30, 255})
		}
	}

	if ui.game.NextPiece != nil {
		shape := ui.game.NextPiece.Shape
		pieceColor := TetrominoColors[ui.game.NextPiece.Type]
		c := color.RGBA{pieceColor.R, pieceColor.G, pieceColor.B, 255}

		offsetX := (width - len(shape[0])*BlockSize) / 2
		offsetY := (height - len(shape)*BlockSize) / 2

		for y := 0; y < len(shape); y++ {
			for x := 0; x < len(shape[y]); x++ {
				if shape[y][x] {
					px := offsetX + x*BlockSize
					py := offsetY + y*BlockSize
					ui.drawBlock(img, px, py, BlockSize, c)
				}
			}
		}
	}

	return img
}

// drawBlock draws a single block with border
func (ui *GameUI) drawBlock(img *image.RGBA, x, y, size int, c color.RGBA) {
	// Draw block
	for dy := 1; dy < size-1; dy++ {
		for dx := 1; dx < size-1; dx++ {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				img.Set(x+dx, y+dy, c)
			}
		}
	}

	// Draw highlight (top and left)
	highlight := color.RGBA{
		uint8(min(255, int(c.R)+50)),
		uint8(min(255, int(c.G)+50)),
		uint8(min(255, int(c.B)+50)),
		255,
	}
	for i := 1; i < size-1; i++ {
		if x+i < img.Bounds().Dx() && y+1 < img.Bounds().Dy() {
			img.Set(x+i, y+1, highlight)
		}
		if x+1 < img.Bounds().Dx() && y+i < img.Bounds().Dy() {
			img.Set(x+1, y+i, highlight)
		}
	}

	// Draw shadow (bottom and right)
	shadow := color.RGBA{
		uint8(max(0, int(c.R)-50)),
		uint8(max(0, int(c.G)-50)),
		uint8(max(0, int(c.B)-50)),
		255,
	}
	for i := 1; i < size-1; i++ {
		if x+i < img.Bounds().Dx() && y+size-2 < img.Bounds().Dy() {
			img.Set(x+i, y+size-2, shadow)
		}
		if x+size-2 < img.Bounds().Dx() && y+i < img.Bounds().Dy() {
			img.Set(x+size-2, y+i, shadow)
		}
	}
}

// Update updates the UI
func (ui *GameUI) Update() {
	ui.boardCanvas.Refresh()
	ui.nextCanvas.Refresh()
	ui.updateLabelsOnly()
}

// updateLabelsOnly updates only the labels (thread-safe)
func (ui *GameUI) updateLabelsOnly() {
	ui.scoreLabel.SetText("Score: " + strconv.Itoa(ui.game.Score))
	ui.linesLabel.SetText("Lines: " + strconv.Itoa(ui.game.Lines))
	ui.levelLabel.SetText("Level: " + strconv.Itoa(ui.game.Level))

	switch ui.game.State {
	case StateMenu:
		ui.statusLabel.SetText("Press Start to begin")
		ui.startButton.Enable()
		ui.pauseButton.Disable()
		ui.pauseButton.SetText("Pause")
	case StatePlaying:
		ui.statusLabel.SetText("Playing")
		ui.startButton.Disable()
		ui.pauseButton.Enable()
		ui.pauseButton.SetText("Pause")
	case StatePaused:
		ui.statusLabel.SetText("Paused")
		ui.startButton.Disable()
		ui.pauseButton.Enable()
		ui.pauseButton.SetText("Resume")
	case StateGameOver:
		// Check if this is a high score (only once per game)
		if !ui.game.scoreSaved {
			wasHighScore := ui.highScoreMgr.AddScore(ui.game.Score)
			if wasHighScore {
				ui.statusLabel.SetText("Game Over!\nNew High Score!")
			} else {
				ui.statusLabel.SetText("Game Over!")
			}
			ui.highScoreLabel.SetText("High Scores:\n" + ui.formatHighScores())
			ui.game.scoreSaved = true
		}
		ui.startButton.Enable()
		ui.pauseButton.Disable()
	}
}

// Show shows the window
func (ui *GameUI) Show() {
	ui.window.Show()
	ui.windowVisible = true
	ui.updateTrayMenuState()
	// Ensure canvas refreshes after window is shown
	ui.boardCanvas.Refresh()
	ui.nextCanvas.Refresh()
}

// ShowAndRun shows the window and runs the application
func (ui *GameUI) ShowAndRun() {
	// ShowAndRun() shows the window and runs the app loop
	ui.window.ShowAndRun()
}

// onStart handles the start button click
func (ui *GameUI) onStart() {
	ui.game.Start()
	// Force immediate refresh of canvas
	ui.boardCanvas.Refresh()
	ui.nextCanvas.Refresh()
	ui.updateLabelsOnly()
}

// onPause handles the pause button click
func (ui *GameUI) onPause() {
	ui.game.Pause()
	ui.Update()
}

// onResetScores handles the reset scores button click
func (ui *GameUI) onResetScores() {
	ui.showResetConfirmationDialog()
}

// showResetConfirmationDialog shows a confirmation dialog requiring the user to type "KrankyBear"
func (ui *GameUI) showResetConfirmationDialog() {
	// If dialog already exists, bring it to front
	if ui.resetDialog != nil {
		ui.resetDialog.RequestFocus()
		ui.resetDialog.Show()
		return
	}

	// Create a new window for the dialog
	dialogWindow := ui.app.NewWindow("Confirm Reset High Scores")
	dialogWindow.Resize(fyne.NewSize(400, 200))
	dialogWindow.SetFixedSize(true)

	// Track the dialog window
	ui.resetDialog = dialogWindow

	// Reset tracking when window closes
	dialogWindow.SetOnClosed(func() {
		ui.resetDialog = nil
	})

	// Create the confirmation text
	confirmLabel := widget.NewLabel("To reset high scores, please type 'KrankyBear' below:")
	confirmLabel.Wrapping = fyne.TextWrapWord

	// Create text entry for confirmation
	confirmEntry := widget.NewEntry()
	confirmEntry.SetPlaceHolder("Type 'KrankyBear' here...")

	// Create buttons
	confirmButton := widget.NewButton("Reset", func() {
		if confirmEntry.Text == "KrankyBear" {
			ui.highScoreMgr.ResetHighScores()
			ui.highScoreLabel.SetText("High Scores:\n" + ui.formatHighScores())
			dialogWindow.Close()
		} else {
			// Show error message
			confirmEntry.SetText("")
			confirmEntry.SetPlaceHolder("Incorrect! Type 'KrankyBear' to confirm...")
		}
	})
	confirmButton.Disable() // Disable until text matches

	cancelButton := widget.NewButton("Cancel", func() {
		dialogWindow.Close()
	})

	// Enable confirm button only when text matches "KrankyBear"
	confirmEntry.OnChanged = func(text string) {
		if text == "KrankyBear" {
			confirmButton.Enable()
		} else {
			confirmButton.Disable()
		}
	}

	// Handle Enter key in entry field
	confirmEntry.OnSubmitted = func(text string) {
		if text == "KrankyBear" {
			confirmButton.OnTapped()
		}
	}

	// Layout
	buttonContainer := container.NewHBox(
		cancelButton,
		confirmButton,
	)

	content := container.NewVBox(
		confirmLabel,
		confirmEntry,
		widget.NewSeparator(),
		buttonContainer,
	)

	dialogWindow.SetContent(content)
	dialogWindow.Show()

	// Position dialog relative to main window (after showing)
	ui.positionDialogRelativeToMain(dialogWindow)

	// Focus the entry field
	confirmEntry.FocusGained()
}

// onToggleSound handles the sound toggle button click
func (ui *GameUI) onToggleSound() {
	ui.soundMgr.SetEnabled(!ui.soundMgr.IsEnabled())
	ui.updateSoundButton()
}

// updateSoundButton updates the sound button text
func (ui *GameUI) updateSoundButton() {
	if ui.soundMgr.IsEnabled() {
		ui.soundButton.SetText("🔊 Sound On")
	} else {
		ui.soundButton.SetText("🔇 Sound Off")
	}
}

// handleSoundEvent handles sound events from the game
func (ui *GameUI) handleSoundEvent(event SoundEvent) {
	if ui.soundMgr == nil {
		return
	}

	switch event {
	case SoundEventPlace:
		ui.soundMgr.PlayPlace()
	case SoundEventLineClear:
		ui.soundMgr.PlayLineClear()
	case SoundEventMultiLineClear:
		ui.soundMgr.PlayMultiLineClear()
	case SoundEventGameOver:
		ui.soundMgr.PlayGameOver()
	}
}

// formatHighScores formats the high scores for display
func (ui *GameUI) formatHighScores() string {
	scores := ui.highScoreMgr.GetHighScores()
	if len(scores) == 0 {
		return "  (none yet)"
	}

	var lines []string
	for i, score := range scores {
		lines = append(lines, "  "+strconv.Itoa(i+1)+". "+strconv.Itoa(score))
	}
	return strings.Join(lines, "\n")
}

// bossKey handles the boss key (F12) - pauses the game and hides the window
func (ui *GameUI) bossKey() {
	// Pause the game if it's currently playing
	if ui.game.State == StatePlaying {
		ui.game.Pause()
		ui.Update()
	}

	// Hide the window
	ui.hideWindow()
}

// HandleKey handles keyboard input
func (ui *GameUI) HandleKey(key fyne.KeyName) {
	switch ui.game.State {
	case StatePlaying:
		switch key {
		case fyne.KeyLeft:
			if ui.game.MoveLeft() {
				ui.Update()
			}
		case fyne.KeyRight:
			if ui.game.MoveRight() {
				ui.Update()
			}
		case fyne.KeyDown:
			if ui.game.MoveDown() {
				ui.Update()
			}
		case fyne.KeyUp:
			if ui.game.RotatePiece() {
				ui.Update()
			}
		case fyne.KeySpace:
			ui.game.HardDrop()
			ui.Update()
		}
	case StatePaused:
		if key == fyne.KeyP || key == fyne.KeySpace {
			ui.game.Pause()
			ui.Update()
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// positionDialogRelativeToMain positions a dialog window relative to the main window
// Uses the same pattern as KrankyBearFileMover: show first, then center
// This ensures the dialog appears on the same display as the main window
func (ui *GameUI) positionDialogRelativeToMain(dialogWindow fyne.Window) {
	// Show the dialog first (it will appear on the same display as the main window)
	// then center it on that display
	dialogWindow.Show()
	dialogWindow.CenterOnScreen()
}

// showAboutDialog shows the About dialog with image on left and text on right
func (ui *GameUI) showAboutDialog() {
	// If dialog already exists, bring it to front
	if ui.aboutDialog != nil {
		ui.aboutDialog.RequestFocus()
		ui.aboutDialog.Show()
		return
	}

	// Create a new window for the dialog
	dialogWindow := ui.app.NewWindow("About")
	dialogWindow.Resize(fyne.NewSize(500, 300))
	dialogWindow.SetFixedSize(true)

	// Track the dialog window
	ui.aboutDialog = dialogWindow

	// Reset tracking when window closes
	dialogWindow.SetOnClosed(func() {
		ui.aboutDialog = nil
	})

	// Create image based on version status
	// Viking Helmet: current version or update available
	// Hard Hat: running newer version than released
	var iconResource fyne.Resource
	if ui.versionStatus == "newer" {
		iconResource = resourceKrankyBearHardHatPng
	} else {
		// "current", "update", or "unknown" - show Viking Helmet
		iconResource = resourceKrankyBearVikingHelmetPng
	}
	iconImage := canvas.NewImageFromResource(iconResource)
	iconImage.FillMode = canvas.ImageFillContain
	iconImage.SetMinSize(fyne.NewSize(150, 150))

	// Create GitHub URL
	githubLink, err := url.Parse("https://github.com/amarillier/KrankyBearTetris")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	githubHyperlink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearTetris", githubLink)
	githubHyperlink.Alignment = fyne.TextAlignLeading

	// Create text content
	aboutText := widget.NewRichTextFromMarkdown(`# KrankyBear Tetris

**Version:** ` + appVersion + `

**Author:** ` + appAuthor + `

**Copyright:** ` + appCopyright + `

A classic Tetris game built with Go and Fyne.

Enjoy the game!`)

	aboutText.Wrapping = fyne.TextWrapWord

	aboutText.Wrapping = fyne.TextWrapWord

	// Layout: image on left, text on right (left-justified)
	textContainer := container.NewVBox(
		container.NewPadded(aboutText),
		githubHyperlink,
	)
	content := container.NewBorder(
		nil,
		container.NewCenter(widget.NewButton("Close", func() {
			dialogWindow.Close()
		})),
		nil,
		nil,
		container.NewHBox(
			container.NewPadded(iconImage),
			textContainer,
		),
	)

	dialogWindow.SetContent(content)
	dialogWindow.Show()

	// Position dialog relative to main window (after showing)
	ui.positionDialogRelativeToMain(dialogWindow)
}

// showUpdateDialog shows the Update Check dialog with version status and appropriate image
// Matches the implementation pattern from KrankyBearTailer
func (ui *GameUI) showUpdateDialog() {
	// Check if window already exists
	if ui.updateDialog != nil {
		ui.updateDialog.RequestFocus()
		return
	}

	// Run update check in goroutine to avoid blocking UI
	go func() {
		updtmsg, _ := ui.updateChecker()
		fyne.Do(func() {
			ui.updateAlert(updtmsg)
		})
	}()
}

// updateAlert displays the update check result in a dialog window
// Matches the implementation pattern from KrankyBearTailer
func (ui *GameUI) updateAlert(updtmsg string) {
	// Check if window already exists
	if ui.updateDialog != nil {
		ui.updateDialog.RequestFocus()
		return
	}

	// Create release link
	releaselink, rerr := url.Parse("https://github.com/amarillier/KrankyBearTetris/releases/latest")
	if rerr != nil {
		fyne.LogError("Could not parse URL", rerr)
	}
	myreleaselink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearTetris/releases/latest", releaselink)
	myreleaselink.Alignment = fyne.TextAlignLeading

	// Create release notes link
	releasenoteslink, rnerr := url.Parse("https://github.com/amarillier/KrankyBearTetris/blob/allanm/ReleaseNotes.txt")
	if rnerr != nil {
		fyne.LogError("Could not parse URL", rnerr)
	}
	myreleasenoteslink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearTetris/blob/allanm/ReleaseNotes.txt", releasenoteslink)
	myreleasenoteslink.Alignment = fyne.TextAlignLeading

	// Create image based on update message
	// Hard Hat: running newer version than released
	// Viking Helmet: current version or update available
	var kbimg *canvas.Image
	if strings.Contains(updtmsg, "newer version") {
		kbimg = canvas.NewImageFromResource(resourceKrankyBearHardHatPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	} else if strings.Contains(updtmsg, "running the latest") {
		kbimg = canvas.NewImageFromResource(resourceKrankyBearVikingHelmetPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	} else {
		// For errors or unknown status, show Viking Helmet
		kbimg = canvas.NewImageFromResource(resourceKrankyBearVikingHelmetPng)
		kbimg.FillMode = canvas.ImageFillOriginal
	}

	// Create text label with update message
	text := widget.NewLabel(updtmsg)
	text.Wrapping = fyne.TextWrapWord

	// Create content: image, text, release link, release notes link
	content := container.NewVBox(kbimg, text, myreleaselink, myreleasenoteslink)

	// Create window
	ui.updateDialog = ui.app.NewWindow(appName + ": Update Check")
	ui.updateDialog.SetIcon(resourceKrankyBearVikingHelmetPng)
	ui.updateDialog.Resize(fyne.NewSize(500, 300))
	ui.updateDialog.SetContent(content)
	ui.updateDialog.SetCloseIntercept(func() {
		ui.updateDialog.Close()
		ui.updateDialog = nil
	})
	ui.updateDialog.Show()
}

// showHelpDialog shows the Help dialog with image on left and text on right
func (ui *GameUI) showHelpDialog() {
	// If dialog already exists, bring it to front
	if ui.helpDialog != nil {
		ui.helpDialog.RequestFocus()
		ui.helpDialog.Show()
		return
	}

	// Create a new window for the dialog
	dialogWindow := ui.app.NewWindow("Help")
	dialogWindow.Resize(fyne.NewSize(500, 400))
	dialogWindow.SetFixedSize(true)

	// Track the dialog window
	ui.helpDialog = dialogWindow

	// Reset tracking when window closes
	dialogWindow.SetOnClosed(func() {
		ui.helpDialog = nil
	})

	// Create image based on version status
	// Viking Helmet: current version or update available
	// Hard Hat: running newer version than released
	var iconResource fyne.Resource
	if ui.versionStatus == "newer" {
		iconResource = resourceKrankyBearHardHatPng
	} else {
		// "current", "update", or "unknown" - show Viking Helmet
		iconResource = resourceKrankyBearVikingHelmetPng
	}
	iconImage := canvas.NewImageFromResource(iconResource)
	iconImage.FillMode = canvas.ImageFillContain
	iconImage.SetMinSize(fyne.NewSize(150, 150))

	// Create GitHub URL
	githubLink, err := url.Parse("https://github.com/amarillier/KrankyBearTetris")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	githubHyperlink := widget.NewHyperlink("https://github.com/amarillier/KrankyBearTetris", githubLink)
	githubHyperlink.Alignment = fyne.TextAlignLeading

	// Create text content
	helpText := widget.NewRichTextFromMarkdown(`# How to Play

## Controls

**Arrow Keys:**
- **← →** Move piece left/right
- **↓** Soft drop (move piece down faster)
- **↑** Rotate piece clockwise

**Space Bar:** Hard drop (instantly drop piece to bottom)

**F12:** Boss key - pause the game and hide the window (quickly hide the game)

## Gameplay

- Clear lines by filling them completely
- Clearing multiple lines at once gives more points
- Game speed increases as you level up
- Try to achieve the highest score!

## Menu Options

- **About:** View application information
- **Help:** Show this help dialog
- **Show:** Show the game window
- **Hide:** Hide the game window (access via system tray)

## System Tray

The application runs in the system tray. Right-click the tray icon to access all menu options.`)

	helpText.Wrapping = fyne.TextWrapWord

	// Layout: image on left, text on right (left-justified)
	textContainer := container.NewVBox(
		container.NewPadded(helpText),
		githubHyperlink,
	)
	content := container.NewBorder(
		nil,
		container.NewCenter(widget.NewButton("Close", func() {
			dialogWindow.Close()
		})),
		nil,
		nil,
		container.NewHBox(
			container.NewPadded(iconImage),
			textContainer,
		),
	)

	dialogWindow.SetContent(content)
	dialogWindow.Show()

	// Position dialog relative to main window (after showing)
	ui.positionDialogRelativeToMain(dialogWindow)
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
