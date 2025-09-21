# TETRIS Game in Go with pterm

A terminal-based TETRIS game implemented in Go using the pterm library for terminal rendering.

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

## Features

- Classic TETRIS gameplay with all 7 tetromino shapes
- Colorful terminal-based graphics
- Score tracking and level progression
- Next piece preview
- Centered display that adapts to terminal size
- Proper terminal cleanup when exiting

## Installation

1. Make sure you have Go installed (version 1.16 or later)
2. Clone or download this repository
3. Navigate to the project directory

## Building

```bash
go build -o tetris cmd/tetris/main.go
```

## Running

```bash
go run cmd/tetris/main.go
```

Or if you built the binary:

```bash
./tetris
```

## Controls

- **← →** : Move left/right
- **↑** : Rotate piece
- **↓** : Soft drop (move down faster)
- **Space** : Hard drop (instantly drop piece)
- **q** : Quit game

## Project Structure

```
pterm-tetris/
├── cmd/
│   └── tetris/
│       └── main.go      # Main application entry point
├── internal/
│   └── game/
│       └── game.go      # Game logic and implementation
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── README.md            # This file
└── LICENSE              # License information
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.