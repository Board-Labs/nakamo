package main

import (
	"fmt"
	"os"

	"github.com/Board-Labs/nakamo/examples/chess/internal/game"
)

func main() {
	fmt.Println("Terminal Chess Game")
	fmt.Println("------------------")
	fmt.Println("Enter moves in algebraic notation (e.g., 'e2 e4')")
	fmt.Println("Type 'quit' to exit")
	fmt.Println()

	g := game.NewGame()
	g.Run()

	fmt.Println("\nThanks for playing!")
	os.Exit(0)
}
