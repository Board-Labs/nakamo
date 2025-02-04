package game

import (
	"github.com/Board-Labs/nakamo/ecs"
)

// MovementSystem calculates possible moves for pieces
type MovementSystem struct {
	world      *ecs.World
	components *Components
}

func NewMovementSystem(world *ecs.World, components *Components) *MovementSystem {
	return &MovementSystem{
		world:      world,
		components: components,
	}
}

func (s *MovementSystem) Update(dt float64) {
	// Calculate possible moves for each piece
	for entity := range s.components.Pieces.GetAll() {
		piece, _ := s.components.Pieces.Get(entity)
		pos, _ := s.components.Positions.Get(entity)
		moves := calculatePossibleMoves(piece, pos)

		if s.components.Movements.Has(entity) {
			s.components.Movements.Remove(entity)
		}
		s.components.Movements.Add(entity, MovementComponent{
			PossibleMoves: moves,
		})
	}
}

func calculatePossibleMoves(piece PieceComponent, pos PositionComponent) []Position {
	moves := make([]Position, 0)

	switch piece.Type {
	case Pawn:
		direction := 1
		if piece.Color == Black {
			direction = -1
		}
		// Forward move
		newPos := Position{pos.Pos.X, pos.Pos.Y + direction}
		if isValidPosition(newPos) {
			moves = append(moves, newPos)
			// Initial two-square move
			if (piece.Color == White && pos.Pos.Y == 1) || (piece.Color == Black && pos.Pos.Y == 6) {
				moves = append(moves, Position{pos.Pos.X, pos.Pos.Y + 2*direction})
			}
		}
	case Knight:
		knightMoves := []Position{
			{pos.Pos.X + 2, pos.Pos.Y + 1}, {pos.Pos.X + 2, pos.Pos.Y - 1},
			{pos.Pos.X - 2, pos.Pos.Y + 1}, {pos.Pos.X - 2, pos.Pos.Y - 1},
			{pos.Pos.X + 1, pos.Pos.Y + 2}, {pos.Pos.X + 1, pos.Pos.Y - 2},
			{pos.Pos.X - 1, pos.Pos.Y + 2}, {pos.Pos.X - 1, pos.Pos.Y - 2},
		}
		for _, move := range knightMoves {
			if isValidPosition(move) {
				moves = append(moves, move)
			}
		}
	case Bishop:
		// Diagonal moves
		directions := []Position{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
		for _, dir := range directions {
			for i := 1; i < 8; i++ {
				newPos := Position{pos.Pos.X + dir.X*i, pos.Pos.Y + dir.Y*i}
				if isValidPosition(newPos) {
					moves = append(moves, newPos)
				}
			}
		}
	case Rook:
		// Horizontal and vertical moves
		directions := []Position{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, dir := range directions {
			for i := 1; i < 8; i++ {
				newPos := Position{pos.Pos.X + dir.X*i, pos.Pos.Y + dir.Y*i}
				if isValidPosition(newPos) {
					moves = append(moves, newPos)
				}
			}
		}
	case Queen:
		// Combine rook and bishop moves
		directions := []Position{
			{0, 1}, {0, -1}, {1, 0}, {-1, 0},
			{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
		}
		for _, dir := range directions {
			for i := 1; i < 8; i++ {
				newPos := Position{pos.Pos.X + dir.X*i, pos.Pos.Y + dir.Y*i}
				if isValidPosition(newPos) {
					moves = append(moves, newPos)
				}
			}
		}
	case King:
		// One square in any direction
		directions := []Position{
			{0, 1}, {0, -1}, {1, 0}, {-1, 0},
			{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
		}
		for _, dir := range directions {
			newPos := Position{pos.Pos.X + dir.X, pos.Pos.Y + dir.Y}
			if isValidPosition(newPos) {
				moves = append(moves, newPos)
			}
		}
	}

	return moves
}

func isValidPosition(pos Position) bool {
	return pos.X >= 0 && pos.X < 8 && pos.Y >= 0 && pos.Y < 8
}

// RenderSystem handles drawing the chess board and pieces
type RenderSystem struct {
	components *Components
}

func NewRenderSystem(components *Components) *RenderSystem {
	return &RenderSystem{
		components: components,
	}
}

func (s *RenderSystem) Update(dt float64) {
	// Clear screen
	print("\033[H\033[2J")

	// Create empty board
	board := make([][]rune, 8)
	for i := range board {
		board[i] = make([]rune, 8)
		for j := range board[i] {
			if (i+j)%2 == 0 {
				board[i][j] = '.'
			} else {
				board[i][j] = ' '
			}
		}
	}

	// Place pieces on board
	for entity := range s.components.Pieces.GetAll() {
		piece, _ := s.components.Pieces.Get(entity)
		pos, _ := s.components.Positions.Get(entity)
		selected := false
		if s.components.Selectable.Has(entity) {
			sel, _ := s.components.Selectable.Get(entity)
			selected = sel.Selected
		}

		symbol := getPieceSymbol(piece)
		board[pos.Pos.Y][pos.Pos.X] = symbol

		if selected {
			moves, err := s.components.Movements.Get(entity)
			if err == nil {
				for _, move := range moves.PossibleMoves {
					if board[move.Y][move.X] == '.' || board[move.Y][move.X] == ' ' {
						board[move.Y][move.X] = '*'
					}
				}
			}
		}
	}

	// Print board
	println("  a b c d e f g h")
	println("  ---------------")
	for i := 7; i >= 0; i-- {
		print(i+1, "|")
		for j := 0; j < 8; j++ {
			print(string(board[i][j]), " ")
		}
		println("|", i+1)
	}
	println("  ---------------")
	println("  a b c d e f g h")
}

func getPieceSymbol(piece PieceComponent) rune {
	symbols := map[PieceType]rune{
		Pawn:   'p',
		Knight: 'n',
		Bishop: 'b',
		Rook:   'r',
		Queen:  'q',
		King:   'k',
	}

	symbol := symbols[piece.Type]
	if piece.Color == White {
		return toUpper(symbol)
	}
	return symbol
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}
