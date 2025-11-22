package main

import (
	"math/rand"
	"time"
)

// Board dimensions
const (
	BoardWidth  = 10
	BoardHeight = 20
)

// Tetromino types
type TetrominoType int

const (
	I TetrominoType = iota
	O
	T
	S
	Z
	J
	L
)

// Color represents a block color
type Color struct {
	R, G, B uint8
}

// Tetromino colors
var TetrominoColors = map[TetrominoType]Color{
	I: {0, 255, 255}, // Cyan
	O: {255, 255, 0}, // Yellow
	T: {128, 0, 128}, // Purple
	S: {0, 255, 0},   // Green
	Z: {255, 0, 0},   // Red
	J: {0, 0, 255},   // Blue
	L: {255, 165, 0}, // Orange
}

// Tetromino represents a tetris piece
type Tetromino struct {
	Type     TetrominoType
	Shape    [][]bool
	X        int
	Y        int
	Rotation int
}

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

// SoundEvent represents a sound event type
type SoundEvent int

const (
	SoundEventPlace SoundEvent = iota
	SoundEventLineClear
	SoundEventMultiLineClear
	SoundEventGameOver
)

// Game represents the Tetris game
type Game struct {
	Board         [BoardHeight][BoardWidth]*Color
	CurrentPiece  *Tetromino
	NextPiece     *Tetromino
	Score         int
	Lines         int
	Level         int
	State         GameState
	DropTimer     time.Time
	DropInterval  time.Duration
	rng           *rand.Rand
	scoreSaved    bool             // Track if score was already saved
	soundCallback func(SoundEvent) // Callback for sound events
	prevLines     int              // Track previous lines count for sound detection
}

// NewGame creates a new game instance
func NewGame() *Game {
	g := &Game{
		State:        StateMenu,
		DropInterval: 1000 * time.Millisecond, // Start with 1 second drop interval
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return g
}

// SetSoundCallback sets the callback function for sound events
func (g *Game) SetSoundCallback(callback func(SoundEvent)) {
	g.soundCallback = callback
}

// Start starts a new game
func (g *Game) Start() {
	g.Board = [BoardHeight][BoardWidth]*Color{}
	g.Score = 0
	g.Lines = 0
	g.Level = 1
	g.State = StatePlaying
	g.DropInterval = 1000 * time.Millisecond
	g.DropTimer = time.Now()
	g.scoreSaved = false // Reset score saved flag
	// Initialize pieces - spawnPiece will handle positioning
	g.CurrentPiece = nil
	g.NextPiece = g.newRandomPiece()
	g.spawnPiece()
}

// spawnPiece places the current piece at the top center of the board
func (g *Game) spawnPiece() {
	// Advance to next piece if CurrentPiece is nil or if we're spawning a new piece
	if g.CurrentPiece == nil {
		// First piece - use NextPiece
		g.CurrentPiece = g.NextPiece
		g.NextPiece = g.newRandomPiece()
	} else {
		// Piece was just placed - advance to next piece
		g.CurrentPiece = g.NextPiece
		g.NextPiece = g.newRandomPiece()
	}

	// Reset position and rotation for the new piece
	g.CurrentPiece.X = BoardWidth/2 - len(g.CurrentPiece.Shape[0])/2
	g.CurrentPiece.Y = 0
	g.CurrentPiece.Rotation = 0

	// Reset drop timer for new piece
	g.DropTimer = time.Now()

	// Check for game over
	if g.checkCollision(g.CurrentPiece) {
		g.State = StateGameOver
		if g.soundCallback != nil {
			g.soundCallback(SoundEventGameOver)
		}
	}
}

// newRandomPiece creates a random tetromino
func (g *Game) newRandomPiece() *Tetromino {
	types := []TetrominoType{I, O, T, S, Z, J, L}
	t := types[g.rng.Intn(len(types))]
	return NewTetromino(t)
}

// NewTetromino creates a new tetromino of the given type
func NewTetromino(t TetrominoType) *Tetromino {
	piece := &Tetromino{
		Type:     t,
		Rotation: 0,
	}

	switch t {
	case I:
		piece.Shape = [][]bool{
			{false, false, false, false},
			{true, true, true, true},
			{false, false, false, false},
			{false, false, false, false},
		}
	case O:
		piece.Shape = [][]bool{
			{true, true},
			{true, true},
		}
	case T:
		piece.Shape = [][]bool{
			{false, true, false},
			{true, true, true},
			{false, false, false},
		}
	case S:
		piece.Shape = [][]bool{
			{false, true, true},
			{true, true, false},
			{false, false, false},
		}
	case Z:
		piece.Shape = [][]bool{
			{true, true, false},
			{false, true, true},
			{false, false, false},
		}
	case J:
		piece.Shape = [][]bool{
			{true, false, false},
			{true, true, true},
			{false, false, false},
		}
	case L:
		piece.Shape = [][]bool{
			{false, false, true},
			{true, true, true},
			{false, false, false},
		}
	}

	return piece
}

// Rotate rotates the tetromino shape
func (t *Tetromino) Rotate() [][]bool {
	size := len(t.Shape)
	rotated := make([][]bool, size)
	for i := range rotated {
		rotated[i] = make([]bool, size)
	}

	// Rotate 90 degrees clockwise
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			rotated[j][size-1-i] = t.Shape[i][j]
		}
	}

	return rotated
}

// GetRotatedShape returns the shape rotated by the current rotation
func (t *Tetromino) GetRotatedShape() [][]bool {
	shape := t.Shape
	for i := 0; i < t.Rotation%4; i++ {
		size := len(shape)
		rotated := make([][]bool, size)
		for j := range rotated {
			rotated[j] = make([]bool, size)
		}
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				rotated[k][size-1-j] = shape[j][k]
			}
		}
		shape = rotated
	}
	return shape
}

// checkCollision checks if the piece collides with the board boundaries or placed blocks
func (g *Game) checkCollision(piece *Tetromino) bool {
	shape := piece.GetRotatedShape()

	for y := 0; y < len(shape); y++ {
		for x := 0; x < len(shape[y]); x++ {
			if !shape[y][x] {
				continue
			}

			boardX := piece.X + x
			boardY := piece.Y + y

			// Check boundaries
			if boardX < 0 || boardX >= BoardWidth || boardY >= BoardHeight {
				return true
			}

			// Check collision with placed blocks (but allow negative Y for spawning)
			if boardY >= 0 && g.Board[boardY][boardX] != nil {
				return true
			}
		}
	}

	return false
}

// MoveLeft moves the current piece left
func (g *Game) MoveLeft() bool {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return false
	}

	g.CurrentPiece.X--
	if g.checkCollision(g.CurrentPiece) {
		g.CurrentPiece.X++
		return false
	}
	return true
}

// MoveRight moves the current piece right
func (g *Game) MoveRight() bool {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return false
	}

	g.CurrentPiece.X++
	if g.checkCollision(g.CurrentPiece) {
		g.CurrentPiece.X--
		return false
	}
	return true
}

// MoveDown moves the current piece down
func (g *Game) MoveDown() bool {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return false
	}

	g.CurrentPiece.Y++
	if g.checkCollision(g.CurrentPiece) {
		g.CurrentPiece.Y--
		g.placePiece()
		return false
	}
	return true
}

// RotatePiece rotates the current piece
func (g *Game) RotatePiece() bool {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return false
	}

	oldRotation := g.CurrentPiece.Rotation
	g.CurrentPiece.Rotation = (g.CurrentPiece.Rotation + 1) % 4

	if g.checkCollision(g.CurrentPiece) {
		// Try wall kicks
		originalX := g.CurrentPiece.X

		// Try moving left
		g.CurrentPiece.X--
		if !g.checkCollision(g.CurrentPiece) {
			return true
		}
		g.CurrentPiece.X = originalX

		// Try moving right
		g.CurrentPiece.X++
		if !g.checkCollision(g.CurrentPiece) {
			return true
		}
		g.CurrentPiece.X = originalX

		// Revert rotation if all attempts fail
		g.CurrentPiece.Rotation = oldRotation
		return false
	}

	return true
}

// HardDrop drops the piece to the bottom instantly
func (g *Game) HardDrop() {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return
	}

	for g.MoveDown() {
		g.Score += 2 // Bonus points for hard drop
	}
}

// placePiece places the current piece on the board
func (g *Game) placePiece() {
	if g.CurrentPiece == nil {
		return
	}

	shape := g.CurrentPiece.GetRotatedShape()
	color := TetrominoColors[g.CurrentPiece.Type]

	for y := 0; y < len(shape); y++ {
		for x := 0; x < len(shape[y]); x++ {
			if !shape[y][x] {
				continue
			}

			boardX := g.CurrentPiece.X + x
			boardY := g.CurrentPiece.Y + y

			if boardY >= 0 && boardY < BoardHeight && boardX >= 0 && boardX < BoardWidth {
				g.Board[boardY][boardX] = &color
			}
		}
	}

	// Clear lines and update score
	linesCleared := g.clearLines()
	g.updateScore(linesCleared)

	// Trigger sound events
	if g.soundCallback != nil {
		g.soundCallback(SoundEventPlace)
		if linesCleared > 0 {
			if linesCleared >= 2 {
				g.soundCallback(SoundEventMultiLineClear)
			} else {
				g.soundCallback(SoundEventLineClear)
			}
		}
	}

	// Spawn next piece
	g.spawnPiece()
}

// clearLines clears completed lines and returns the number cleared
func (g *Game) clearLines() int {
	linesCleared := 0

	for y := BoardHeight - 1; y >= 0; y-- {
		full := true
		for x := 0; x < BoardWidth; x++ {
			if g.Board[y][x] == nil {
				full = false
				break
			}
		}

		if full {
			// Remove the line
			for y2 := y; y2 > 0; y2-- {
				g.Board[y2] = g.Board[y2-1]
			}
			// Clear top line
			for x := 0; x < BoardWidth; x++ {
				g.Board[0][x] = nil
			}
			linesCleared++
			y++ // Check the same line again (it's now the line above)
		}
	}

	return linesCleared
}

// updateScore updates the score based on lines cleared
func (g *Game) updateScore(linesCleared int) {
	if linesCleared == 0 {
		return
	}

	// Scoring: 1 line = 100, 2 lines = 300, 3 lines = 500, 4 lines = 800
	points := []int{0, 100, 300, 500, 800}
	if linesCleared < len(points) {
		g.Score += points[linesCleared] * (g.Level + 1)
	} else {
		g.Score += 800 * (g.Level + 1)
	}

	g.Lines += linesCleared
	// Level up every 10 lines
	newLevel := g.Lines/10 + 1
	if newLevel > g.Level {
		g.Level = newLevel
		// Speed up: reduce drop interval by 50ms per level, minimum 100ms
		newInterval := 1000*time.Millisecond - time.Duration(g.Level-1)*50*time.Millisecond
		if newInterval < 100*time.Millisecond {
			newInterval = 100 * time.Millisecond
		}
		g.DropInterval = newInterval
	}
}

// Update updates the game state (called every frame)
func (g *Game) Update() {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return
	}

	// Auto-drop
	if time.Since(g.DropTimer) >= g.DropInterval {
		g.MoveDown()
		g.DropTimer = time.Now()
	}
}

// Pause toggles pause state
func (g *Game) Pause() {
	if g.State == StatePlaying {
		g.State = StatePaused
	} else if g.State == StatePaused {
		g.State = StatePlaying
		g.DropTimer = time.Now()
	}
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
