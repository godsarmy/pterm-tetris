package main

import (
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	"pterm-tetris/internal/game"
)

func main() {
	// Seed random number generator
	game.SeedRandom()

	// Show welcome screen
	game.ShowWelcomeScreen()

	// Wait for key press to start
	game.WaitForStart()

	// Clear screen again for game
	print("\033[H\033[2J")

	// Create game
	g := game.NewGame()

	// Create area for rendering
	area, err := pterm.DefaultArea.Start()
	if err != nil {
		panic(err)
	}
	defer area.Stop()

	// Channel to signal game exit
	exitChan := make(chan bool, 1) // Buffered channel to prevent blocking

	// Channel to signal keyboard listener exit
	keyboardDone := make(chan bool, 1) // Buffered channel

	// Start keyboard listener in a separate goroutine
	go func() {
		defer func() {
			keyboardDone <- true
		}()

		err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
			if key.Code == keys.RuneKey && (key.String() == "q" || key.String() == "Q") {
				exitChan <- true
				return true, nil // Quit game
			}

			g.HandleInput(key)
			return false, nil
		})

		if err != nil {
			pterm.Error.Printfln("Failed to start keyboard listener: %v", err)
		}
	}()

	// Game loop
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	// Flag to track if we need to reset terminal
	needReset := false
	gameFinished := false

	for {
		select {
		case <-exitChan:
			needReset = true
			gameFinished = true
		case <-ticker.C:
			if !gameFinished {
				g.Update()
				g.Draw(area)

				if g.GameOver {
					// Keep the game display visible after game over
					time.Sleep(3 * time.Second)
					needReset = true
					gameFinished = true
				}
			}
		default:
			if gameFinished {
				goto cleanup
			}
			time.Sleep(10 * time.Millisecond)
		}

		if gameFinished {
			goto cleanup
		}
	}

cleanup:
	// Stop the area printer
	area.Stop()

	// Wait for keyboard listener to finish with a timeout
	select {
	case <-keyboardDone:
		// Keyboard listener finished normally
	case <-time.After(100 * time.Millisecond):
		// Timeout - keyboard listener didn't finish, but we need to continue
	}

	// Reset terminal if needed
	if needReset {
		game.CleanupTerminal()
	}
}
