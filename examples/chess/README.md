# Terminal Chess Game

A terminal-based chess game implemented using the Entity Component System (ECS) framework.

## Features

- Full chess game implementation
- Terminal-based UI with ASCII graphics
- Move validation
- Two-player gameplay
- Standard algebraic notation for moves

## How to Run

1. Make sure you have Go installed on your system
2. Navigate to the chess game directory:
   ```bash
   cd examples/chess/cmd/chess
   ```
3. Run the game:
   ```bash
   go run main.go
   ```

## How to Play

1. The game uses standard algebraic notation for moves
2. Enter moves in the format: `e2 e4` (source square to destination square)
3. White pieces are represented by uppercase letters (P, N, B, R, Q, K)
4. Black pieces are represented by lowercase letters (p, n, b, r, q, k)
5. Type 'quit' to exit the game

## Project Structure

```
chess/
├── cmd/
│   └── chess/
│       └── main.go         # Entry point
├── internal/
│   └── game/
│       ├── components.go   # Game components
│       ├── systems.go      # Game systems
│       └── game.go         # Game logic
└── README.md
```

## Implementation Details

The game is built using an Entity Component System (ECS) architecture:

- **Components:**
  - Position: Stores piece location
  - Piece: Stores piece type and color
  - Movement: Stores possible moves
  - Selectable: Handles piece selection

- **Systems:**
  - MovementSystem: Calculates valid moves
  - RenderSystem: Handles board display

The ECS architecture provides a clean separation of concerns and makes the code modular and extensible.