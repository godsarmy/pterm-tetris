package game

import (
	"math/rand"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	"pterm-tetris/internal/version"
)

// Define the tetromino shapes
type Tetromino struct {
	Shape    [][]int
	Color    pterm.Color
	X, Y     int
	Rotation int
}

// Define the game board
type Board struct {
	Width, Height int
	Grid          [][]int
}

// Define the game state
type Game struct {
	Board        Board
	Current      Tetromino
	Next         Tetromino
	Score        int
	Level        int
	Lines        int
	GameOver     bool
	DropTime     time.Time
	DropSpeed    time.Duration
	GhostEnabled bool
}

// Tetromino shapes
var (
	I = Tetromino{
		Shape: [][]int{
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		Color: pterm.FgCyan,
	}

	O = Tetromino{
		Shape: [][]int{
			{1, 1},
			{1, 1},
		},
		Color: pterm.FgYellow,
	}

	T = Tetromino{
		Shape: [][]int{
			{0, 1, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		Color: pterm.FgMagenta,
	}

	S = Tetromino{
		Shape: [][]int{
			{0, 1, 1},
			{1, 1, 0},
			{0, 0, 0},
		},
		Color: pterm.FgGreen,
	}

	Z = Tetromino{
		Shape: [][]int{
			{1, 1, 0},
			{0, 1, 1},
			{0, 0, 0},
		},
		Color: pterm.FgRed,
	}

	J = Tetromino{
		Shape: [][]int{
			{1, 0, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		Color: pterm.FgBlue,
	}

	L = Tetromino{
		Shape: [][]int{
			{0, 0, 1},
			{1, 1, 1},
			{0, 0, 0},
		},
		Color: pterm.FgLightYellow,
	}
)

// All tetrominoes
var tetrominoes = []Tetromino{I, O, T, S, Z, J, L}

// Initialize the game board
func NewBoard(width, height int) Board {
	grid := make([][]int, height)
	for i := range grid {
		grid[i] = make([]int, width)
	}
	return Board{Width: width, Height: height, Grid: grid}
}

// Initialize a new game
func NewGame() *Game {
	board := NewBoard(10, 20)
	game := &Game{
		Board:        board,
		Score:        0,
		Level:        1,
		Lines:        0,
		DropSpeed:    500 * time.Millisecond,
		GhostEnabled: true,
	}
	game.Current = game.newTetromino()
	game.Next = game.newTetromino()
	game.DropTime = time.Now().Add(game.DropSpeed)
	return game
}

// Create a new random tetromino
func (g *Game) newTetromino() Tetromino {
	tet := tetrominoes[rand.Intn(len(tetrominoes))]
	tet.X = g.Board.Width/2 - len(tet.Shape[0])/2
	tet.Y = 0
	return tet
}

// Rotate a tetromino
func (t *Tetromino) Rotate() {
	// Create a new rotated shape
	size := len(t.Shape)
	rotated := make([][]int, size)
	for i := range rotated {
		rotated[i] = make([]int, size)
	}

	// Perform rotation
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			rotated[x][size-1-y] = t.Shape[y][x]
		}
	}

	t.Shape = rotated
}

// Check for collision
func (g *Game) CheckCollision(t Tetromino, dx, dy int) bool {
	for y := 0; y < len(t.Shape); y++ {
		for x := 0; x < len(t.Shape[y]); x++ {
			if t.Shape[y][x] != 0 {
				newX := t.X + x + dx
				newY := t.Y + y + dy

				// Check boundaries
				if newX < 0 || newX >= g.Board.Width || newY >= g.Board.Height {
					return true
				}

				// Check if already occupied (but only if not above the board)
				if newY >= 0 && g.Board.Grid[newY][newX] != 0 {
					return true
				}
			}
		}
	}
	return false
}

// Merge tetromino with board
func (g *Game) MergeTetromino() {
	for y := 0; y < len(g.Current.Shape); y++ {
		for x := 0; x < len(g.Current.Shape[y]); x++ {
			if g.Current.Shape[y][x] != 0 {
				boardY := g.Current.Y + y
				boardX := g.Current.X + x
				if boardY >= 0 && boardY < g.Board.Height && boardX >= 0 && boardX < g.Board.Width {
					// Use color value to represent different colors
					g.Board.Grid[boardY][boardX] = int(g.Current.Color)
				}
			}
		}
	}
}

// Clear completed lines
func (g *Game) ClearLines() int {
	linesCleared := 0
	for y := g.Board.Height - 1; y >= 0; y-- {
		full := true
		for x := 0; x < g.Board.Width; x++ {
			if g.Board.Grid[y][x] == 0 {
				full = false
				break
			}
		}
		if full {
			// Move all lines above down
			for y2 := y; y2 > 0; y2-- {
				for x := 0; x < g.Board.Width; x++ {
					g.Board.Grid[y2][x] = g.Board.Grid[y2-1][x]
				}
			}
			// Clear the top line
			for x := 0; x < g.Board.Width; x++ {
				g.Board.Grid[0][x] = 0
			}
			linesCleared++
			y++ // Check the same line again
		}
	}
	return linesCleared
}

// Update game state
func (g *Game) Update() {
	if g.GameOver {
		return
	}

	// Check if it's time to drop
	if time.Now().After(g.DropTime) {
		if !g.CheckCollision(g.Current, 0, 1) {
			g.Current.Y++
		} else {
			// Merge with board
			g.MergeTetromino()

			// Clear lines
			lines := g.ClearLines()
			if lines > 0 {
				g.Lines += lines
				g.Score += lines * 100 * g.Level
				g.Level = g.Lines/10 + 1
				g.DropSpeed = time.Duration(500-50*(g.Level-1)) * time.Millisecond
				if g.DropSpeed < 50*time.Millisecond {
					g.DropSpeed = 50 * time.Millisecond
				}
			}

			// Create new tetromino
			g.Current = g.Next
			g.Next = g.newTetromino()

			// Check for game over
			if g.CheckCollision(g.Current, 0, 0) {
				g.GameOver = true
			}
		}
		g.DropTime = time.Now().Add(g.DropSpeed)
	}
}

// Move current tetromino
func (g *Game) Move(dx, dy int) {
	if !g.CheckCollision(g.Current, dx, dy) {
		g.Current.X += dx
		g.Current.Y += dy
	}
}

// Rotate current tetromino
func (g *Game) Rotate() {
	originalShape := g.Current.Shape
	g.Current.Rotate()
	if g.CheckCollision(g.Current, 0, 0) {
		// Revert if collision
		g.Current.Shape = originalShape
	}
}

// Draw the game
func (g *Game) Draw(area *pterm.AreaPrinter) {
	// Create a copy of the board to draw on
	displayBoard := make([][]int, g.Board.Height)
	for i := range displayBoard {
		displayBoard[i] = make([]int, g.Board.Width)
		copy(displayBoard[i], g.Board.Grid[i])
	}

	// Draw ghost piece (projection of hard drop)
	if g.GhostEnabled && !g.GameOver {
		ghost := g.Current
		for !g.CheckCollision(ghost, 0, 1) {
			ghost.Y++
		}
		// Mark ghost cells with sentinel value -1
		for y := 0; y < len(ghost.Shape); y++ {
			for x := 0; x < len(ghost.Shape[y]); x++ {
				if ghost.Shape[y][x] != 0 {
					boardY := ghost.Y + y
					boardX := ghost.X + x
					if boardY >= 0 && boardY < g.Board.Height && boardX >= 0 && boardX < g.Board.Width {
						if displayBoard[boardY][boardX] == 0 { // don't overwrite existing blocks
							displayBoard[boardY][boardX] = -1
						}
					}
				}
			}
		}
	}

	// Draw current tetromino on the board copy
	for y := 0; y < len(g.Current.Shape); y++ {
		for x := 0; x < len(g.Current.Shape[y]); x++ {
			if g.Current.Shape[y][x] != 0 {
				boardY := g.Current.Y + y
				boardX := g.Current.X + x
				if boardY >= 0 && boardY < g.Board.Height && boardX >= 0 && boardX < g.Board.Width {
					displayBoard[boardY][boardX] = int(g.Current.Color)
				}
			}
		}
	}

	// Calculate padding for centering
	terminalWidth := 80  // Default width if we can't get actual terminal size
	terminalHeight := 24 // Default height if we can't get actual terminal size

	// Try to get actual terminal size
	if width, height, err := pterm.GetTerminalSize(); err == nil {
		terminalWidth = width
		terminalHeight = height
	}

	// Prepare game board content
	var boardLines []string
	boardLines = append(boardLines, "┌────────────────────┐")
	for y := 0; y < g.Board.Height; y++ {
		line := "│"
		for x := 0; x < g.Board.Width; x++ {
			switch displayBoard[y][x] {
			case 0: // Empty
				line += "  "
			case -1: // Ghost piece
				line += pterm.FgGray.Sprint("░░")
			default: // Placed piece or current piece
				// Convert back to color
				color := pterm.Color(displayBoard[y][x])
				line += color.Sprint("██")
			}
		}
		line += "│"
		boardLines = append(boardLines, line)
	}
	boardLines = append(boardLines, "└────────────────────┘")

	// Prepare info panel content
	var infoLines []string
	infoLines = append(infoLines, pterm.FgLightCyan.Sprint("TETRIS"))
	infoLines = append(infoLines, pterm.FgLightWhite.Sprint("v"+version.Version))
	infoLines = append(infoLines, "")
	infoLines = append(infoLines, pterm.Sprintf("Score: %d", g.Score))
	infoLines = append(infoLines, pterm.Sprintf("Level: %d", g.Level))
	infoLines = append(infoLines, pterm.Sprintf("Lines: %d", g.Lines))
	infoLines = append(infoLines, "")
	infoLines = append(infoLines, "Next:")

	// Draw next piece preview
	nextSize := len(g.Next.Shape)
	for y := 0; y < nextSize; y++ {
		line := ""
		for x := 0; x < nextSize; x++ {
			if y < len(g.Next.Shape) && x < len(g.Next.Shape[y]) && g.Next.Shape[y][x] != 0 {
				line += g.Next.Color.Sprint("██")
			} else {
				line += "  "
			}
		}
		infoLines = append(infoLines, line)
	}

	infoLines = append(infoLines, "")
	infoLines = append(infoLines, "Controls:")
	infoLines = append(infoLines, "← → : Move")
	infoLines = append(infoLines, "↑   : Rotate")
	infoLines = append(infoLines, "↓   : Soft Drop")
	infoLines = append(infoLines, "Space : Hard Drop")
	infoLines = append(infoLines, "g   : Toggle Ghost")
	infoLines = append(infoLines, "q   : Quit")

	if g.GameOver {
		infoLines = append(infoLines, "")
		infoLines = append(infoLines, pterm.FgRed.Sprint("GAME OVER!"))
		infoLines = append(infoLines, pterm.FgRed.Sprint("Press 'q' to quit."))
	}

	// Calculate layout
	boardWidth := 24 // 2 borders + 2*10 blocks + 2 spaces
	infoWidth := 20
	totalContentWidth := boardWidth + infoWidth + 2 // +2 for spacing

	horizontalPadding := (terminalWidth - totalContentWidth) / 2
	verticalPadding := (terminalHeight - len(boardLines) - 2) / 2 // -2 for title lines

	// Ensure padding is not negative
	if horizontalPadding < 0 {
		horizontalPadding = 0
	}
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	// Create horizontal padding
	hPadding := ""
	for i := 0; i < horizontalPadding; i++ {
		hPadding += " "
	}

	// Draw the content
	content := "\n"

	// Add vertical padding
	for i := 0; i < verticalPadding/2; i++ {
		content += "\n"
	}

	// Draw the game area with two columns
	for i := 0; i < len(boardLines) || i < len(infoLines); i++ {
		content += hPadding // Add horizontal padding

		// Left column (game board)
		if i < len(boardLines) {
			content += boardLines[i]
		} else {
			// Fill with empty space to match board height
			content += "                        " // 24 spaces
		}

		// Space between columns
		content += "  "

		// Right column (info panel)
		if i < len(infoLines) {
			content += infoLines[i]
		}

		content += "\n"
	}

	// Add remaining vertical padding
	for i := 0; i < verticalPadding/2; i++ {
		content += "\n"
	}

	area.Update(content)
}

// Handle keyboard input
func (g *Game) HandleInput(key keys.Key) {
	switch key.Code {
	case keys.RuneKey:
		switch key.String() {
		case "a", "A":
			g.Move(-1, 0) // Move left
		case "d", "D":
			g.Move(1, 0) // Move right
		case "s", "S":
			g.Move(0, 1) // Move down
		case "w", "W":
			g.Rotate() // Rotate
		case "g", "G":
			g.GhostEnabled = !g.GhostEnabled // Toggle ghost piece
		case " ":
			// Hard drop
			for !g.CheckCollision(g.Current, 0, 1) {
				g.Current.Y++
			}
		case "q", "Q":
			g.GameOver = true // Quit game
		}
	case keys.Enter:
		g.Rotate() // Rotate
	case keys.Left:
		g.Move(-1, 0) // Move left
	case keys.Right:
		g.Move(1, 0) // Move right
	case keys.Down:
		g.Move(0, 1) // Move down
	case keys.Up:
		g.Rotate() // Rotate
	case keys.Space:
		// Hard drop
		for !g.CheckCollision(g.Current, 0, 1) {
			g.Current.Y++
		}
	}
}

// SeedRandom seeds the random number generator
func SeedRandom() {
	rand.Seed(time.Now().UnixNano())
}

// ShowWelcomeScreen shows the welcome screen with ASCII art
func ShowWelcomeScreen() {
	// Clear screen
	print("\033[H\033[2J")

	// Calculate padding for centering welcome message
	terminalWidth := 80  // Default width if we can't get actual terminal size
	terminalHeight := 24 // Default height if we can't get actual terminal size

	// Try to get actual terminal size
	if width, height, err := pterm.GetTerminalSize(); err == nil {
		terminalWidth = width
		terminalHeight = height
	}

	// Create ASCII art for "TETRIS"
	tetrisArt := []string{
		"████████╗███████╗████████╗██████╗ ██╗███████╗",
		"╚══██╔══╝██╔════╝╚══██╔══╝██╔══██╗██║██╔════╝",
		"   ██║   █████╗     ██║   ██████╔╝██║███████╗",
		"   ██║   ██╔══╝     ██║   ██╔══██╗██║╚════██║",
		"   ██║   ███████╗   ██║   ██║  ██║██║███████║",
		"   ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝",
	}

	// Calculate horizontal padding
	hPadding := ""
	titleWidth := len(tetrisArt[0]) // Width of the first line
	if titleWidth < terminalWidth {
		for i := 0; i < (terminalWidth-titleWidth)/2; i++ {
			hPadding += " "
		}
	}

	// Calculate vertical padding
	vPadding := ""
	artHeight := len(tetrisArt) + 4 // Art lines + version + "Press any key to start" + spacing
	if artHeight < terminalHeight {
		for i := 0; i < (terminalHeight-artHeight)/2; i++ {
			vPadding += "\n"
		}
	}

	// Show centered welcome message with large ASCII art
	content := vPadding

	// Add the ASCII art
	for _, line := range tetrisArt {
		content += hPadding + pterm.FgLightCyan.Sprint(line) + "\n"
	}

	// Add version information
	content += "\n" + hPadding + "              v" + version.Version + "\n"

	// Add spacing and start message
	content += "\n" + hPadding + "        Press any key to start..." + "\n" + vPadding

	// Print the centered content
	print(content)
}

// WaitForStart waits for a key press to start the game
func WaitForStart() {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		return true, nil
	})
}

// CleanupTerminal cleans up the terminal after the game
func CleanupTerminal() {
	print("\033c")     // Reset terminal
	print("\033[?25h") // Show cursor
}
