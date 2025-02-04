package game

import "github.com/Board-Labs/nakamo/ecs"

// Position represents a chess piece's position on the board
type Position struct {
	X, Y int
}

// PositionComponent stores a piece's position
type PositionComponent struct {
	Pos Position
}

// PieceType represents different types of chess pieces
type PieceType int

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

// Color represents the color of a chess piece
type Color int

const (
	White Color = iota
	Black
)

// PieceComponent stores information about a chess piece
type PieceComponent struct {
	Type  PieceType
	Color Color
}

// MovementComponent stores possible moves for a piece
type MovementComponent struct {
	PossibleMoves []Position
}

// SelectableComponent marks entities that can be selected
type SelectableComponent struct {
	Selected bool
}

// Components holds all component storages for the chess game
type Components struct {
	Positions  *ecs.Storage[PositionComponent]
	Pieces     *ecs.Storage[PieceComponent]
	Movements  *ecs.Storage[MovementComponent]
	Selectable *ecs.Storage[SelectableComponent]
}

// NewComponents creates and initializes all component storages
func NewComponents() *Components {
	return &Components{
		Positions:  ecs.NewStorage[PositionComponent](),
		Pieces:     ecs.NewStorage[PieceComponent](),
		Movements:  ecs.NewStorage[MovementComponent](),
		Selectable: ecs.NewStorage[SelectableComponent](),
	}
}

// RegisterComponents registers all component storages with the world
func RegisterComponents(world *ecs.World) *Components {
	components := NewComponents()
	world.RegisterStorage(components.Positions)
	world.RegisterStorage(components.Pieces)
	world.RegisterStorage(components.Movements)
	world.RegisterStorage(components.Selectable)
	return components
}
