package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

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
	debug := flag.Bool("debug", false, "output debug info only in json format")
	timeout := flag.Duration("timeout", time.Duration(0), "timeout for solver")
	file := flag.String("f", "-", "file with map or - for stdin")
	flag.Parse()

	log.SetOutput(os.Stderr)

	var input io.Reader
	if *file == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(*file)
		if err != nil {
			log.Fatalf("File open error: %v", err)
		}
		input = f
		defer f.Close()
	}

	m, err := world.ReadMap(input)
	if err != nil {
		log.Fatalf("Map reading error: %v", err)
	}
	log.Print("Loaded: ", m)

	solver := world.NewSolver(m)
	// TODO: prepare data

	ctx := context.Background()
	if *timeout > time.Duration(0) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}
	log.Print("Starting solution...")
	err = solver.Solve(ctx, *debug)
	if err != nil {
		log.Fatalf("Solving error: %v", err)
	}

	if *debug {
		dbgInfo := solver.GetDebug()
		err := json.NewEncoder(os.Stdout).Encode(dbgInfo)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
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
