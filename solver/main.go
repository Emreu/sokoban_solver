package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/emreu/sokoban_solver/solver/world"
)

func directionToString(dir world.MoveDirection) string {
	switch dir {
	case world.MoveUp:
		return "^"
	case world.MoveRight:
		return ">"
	case world.MoveLeft:
		return "<"
	case world.MoveDown:
		return "V"
	}
	return ""
}

func main() {
	// TODO: read in cmd args

	m, err := world.ReadMap(os.Stdin)
	if err != nil {
		log.Fatalf("Map reading error: %v", err)
	}

	solver := world.NewSolver(m)
	// TODO: prepare data

	// TODO: add timeout to context if required by cmd args
	ctx := context.Background()
	err = solver.Solve(ctx)
	if err != nil {
		log.Fatalf("Solving error: %v", err)
	}

	path, err := solver.GetPath()
	if err != nil {
		log.Fatalf("Path building error: %v", err)
	}
	for _, dir := range path {
		fmt.Println(directionToString(dir))
	}
}
