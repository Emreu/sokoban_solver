package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/emreu/sokoban_solver/solver/server"
	"github.com/emreu/sokoban_solver/solver/world"
)

func main() {

	srv := flag.Bool("server", false, "start server")
	port := flag.Int("port", 3000, "server listen port")

	timeout := flag.Duration("timeout", time.Duration(0), "timeout for solver")
	file := flag.String("f", "-", "file with map or - for stdin")
	flag.Parse()

	if *srv {
		err := server.Run(*port)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

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

	ctx := context.Background()
	if *timeout > time.Duration(0) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}
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
		fmt.Println(dir.String())
	}
}
