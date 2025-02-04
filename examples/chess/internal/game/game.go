package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Board-Labs/nakamo/ecs"
)

// Game represents the chess game
type Game struct {
	world      *ecs.World
	components *Components
	selected   ecs.Entity
}

// NewGame creates a new chess game
func NewGame() *Game {
	world := ecs.NewWorld()
	components := RegisterComponents(world)

	game := &Game{
		world:      world,
		components: components,
		selected:   0,
	}

	// Add systems
	world.AddSystem(NewMovementSystem(world, components))
	world.AddSystem(NewRenderSystem(components))

	// Initialize board
	game.setupBoard()

	return game
}

// setupBoard creates all chess pieces in their initial positions
func (g *Game) setupBoard() {
	// Create white pieces
	g.createPiece(Rook, White, 0, 0)
	g.createPiece(Knight, White, 1, 0)
	g.createPiece(Bishop, White, 2, 0)
	g.createPiece(Queen, White, 3, 0)
	g.createPiece(King, White, 4, 0)
	g.createPiece(Bishop, White, 5, 0)
	g.createPiece(Knight, White, 6, 0)
	g.createPiece(Rook, White, 7, 0)

	for i := 0; i < 8; i++ {
		g.createPiece(Pawn, White, i, 1)
	}

	// Create black pieces
	g.createPiece(Rook, Black, 0, 7)
	g.createPiece(Knight, Black, 1, 7)
	g.createPiece(Bishop, Black, 2, 7)
	g.createPiece(Queen, Black, 3, 7)
	g.createPiece(King, Black, 4, 7)
	g.createPiece(Bishop, Black, 5, 7)
	g.createPiece(Knight, Black, 6, 7)
	g.createPiece(Rook, Black, 7, 7)

	for i := 0; i < 8; i++ {
		g.createPiece(Pawn, Black, i, 6)
	}
}

// createPiece creates a new chess piece entity
func (g *Game) createPiece(pieceType PieceType, color Color, x, y int) ecs.Entity {
	entity := g.world.CreateEntity()

	g.components.Positions.Add(entity, PositionComponent{
		Pos: Position{X: x, Y: y},
	})

	g.components.Pieces.Add(entity, PieceComponent{
		Type:  pieceType,
		Color: color,
	})

	g.components.Selectable.Add(entity, SelectableComponent{
		Selected: false,
	})

	return entity
}

// Run starts the game loop
func (g *Game) Run() {
	reader := bufio.NewReader(os.Stdin)
	currentPlayer := White

	for {
		// Update game state
		g.world.Update(0)

		// Get player input
		fmt.Printf("\nPlayer %s's turn\n", getColorName(currentPlayer))
		fmt.Print("Enter move (e.g. 'e2 e4' or 'quit'): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			break
		}

		// Parse move
		parts := strings.Split(input, " ")
		if len(parts) != 2 {
			fmt.Println("Invalid input format. Use 'e2 e4' format.")
			continue
		}

		from := parsePosition(parts[0])
		to := parsePosition(parts[1])

		if !isValidPosition(from) || !isValidPosition(to) {
			fmt.Println("Invalid position.")
			continue
		}

		// Find piece at source position
		var pieceEntity ecs.Entity
		found := false

		for entity := range g.components.Pieces.GetAll() {
			pos, _ := g.components.Positions.Get(entity)
			if pos.Pos == from {
				piece, _ := g.components.Pieces.Get(entity)
				if piece.Color == currentPlayer {
					pieceEntity = entity
					found = true
					break
				}
			}
		}

		if !found {
			fmt.Println("No piece found at source position.")
			continue
		}

		// Check if move is valid
		moves, _ := g.components.Movements.Get(pieceEntity)
		validMove := false
		for _, move := range moves.PossibleMoves {
			if move == to {
				validMove = true
				break
			}
		}

		if !validMove {
			fmt.Println("Invalid move for this piece.")
			continue
		}

		// Remove any piece at destination
		for entity := range g.components.Pieces.GetAll() {
			pos, _ := g.components.Positions.Get(entity)
			if pos.Pos == to {
				g.world.DestroyEntity(entity)
				break
			}
		}

		// Move piece
		g.components.Positions.Add(pieceEntity, PositionComponent{Pos: to})

		// Switch players
		if currentPlayer == White {
			currentPlayer = Black
		} else {
			currentPlayer = White
		}
	}
}

func getColorName(color Color) string {
	if color == White {
		return "White"
	}
	return "Black"
}

func parsePosition(s string) Position {
	if len(s) != 2 {
		return Position{-1, -1}
	}

	file := int(s[0] - 'a')
	rank := int(s[1] - '1')

	return Position{file, rank}
}
