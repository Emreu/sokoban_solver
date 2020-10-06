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
	log.SetOutput(os.Stderr)

	f, err := os.Open("test_1.txt")
	if err != nil {
		log.Fatalf("File open error: %v", err)
	}
	defer f.Close()

	// m, err := world.ReadMap(os.Stdin)
	m, err := world.ReadMap(f)
	if err != nil {
		log.Fatalf("Map reading error: %v", err)
	}
	log.Print("Loaded: ", m)

	solver := world.NewSolver(m)
	// TODO: prepare data

	// TODO: add timeout to context if required by cmd args
	ctx := context.Background()
	log.Print("Starting solution...")
	err = solver.Solve(ctx)
	if err != nil {
		log.Fatalf("Solving error: %v", err)
	}
	log.Print("Solved!")

	path, err := solver.GetPath()
	if err != nil {
		log.Fatalf("Path building error: %v", err)
	}
	for _, dir := range path {
		fmt.Println(directionToString(dir))
	}
}
