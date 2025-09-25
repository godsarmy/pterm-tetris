package main

import (
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	"pterm-tetris/internal/game"
	_ "pterm-tetris/internal/version"
)

func main() {
	// Seed random number generator
	game.SeedRandom()

	// Show welcome screen
	game.ShowWelcomeScreen()

	// Wait for key press to start
	game.WaitForStart()

	// Cleanup screen again for game
	game.CleanupTerminal()

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
			if key.Code == keys.CtrlC {
				exitChan <- true
				return true, nil // Quit game on Ctrl+C
			}

			// Handle quit confirmation flow while playing
			if g.ConfirmQuit {
				if key.Code == keys.RuneKey {
					switch key.String() {
					case "y", "Y":
						exitChan <- true
						return true, nil // Confirm quit
					case "n", "N":
						g.ConfirmQuit = false
						return false, nil // Cancel quit
					}
				}
				return false, nil // Ignore other keys during quit confirm
			}

			// Immediate quit on 'q' only if game is already over
			if key.Code == keys.RuneKey && (key.String() == "q" || key.String() == "Q") {
				if g.GameOver {
					exitChan <- true
					return true, nil // Quit after game over
				}
				// If not over, trigger confirmation prompt
				g.ConfirmQuit = true
				return false, nil
			}

			if g.GameOver {
				// Allow restart confirmation flow while game over
				if key.Code == keys.RuneKey {
					switch key.String() {
					case "r", "R", "y", "Y", "n", "N":
						g.HandleInput(key)
					}
				}
				return false, nil // Ignore other keys after game over
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
					// Keep the game display visible until explicit quit
					needReset = true
					// Do not set gameFinished; wait for 'q' or Ctrl+C
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
		game.ResetTerminal()
	}
}
