package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/emreu/sokoban_solver/solver/world"
)

type solveResponse struct {
	Solution       []world.MoveDirection `json:"solution,omitempty"`
	DeadZones      world.PosList         `json:"dead_zones,omitempty"`
	HorizontalLock world.PosList         `json:"horizontal_locks,omitempty"`
	VerticalLock   world.PosList         `json:"vertical_locks,omitempty"`
	Metrics        []map[string]int      `json:"metrics,omitempty"`
	States         []*world.Node         `json:"states,omitempty"`
	Error          string                `json:"error,omitempty"`
}

func errorMessage(res http.ResponseWriter, status int, err error) {
	log.Printf("Error processing request: %v", err)
	res.WriteHeader(status)
	fmt.Fprintf(res, "Can't process request: %v", err)
}

func solveHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("Serving <solve> request from %s...", req.RemoteAddr)
	if req.Method != http.MethodPost {
		errorMessage(res, http.StatusMethodNotAllowed, fmt.Errorf("non-POST method"))
		return
	}

	// get level from request
	m, err := world.ReadMap(req.Body)
	if err != nil {
		errorMessage(res, http.StatusBadRequest, fmt.Errorf("map reading error: %v", err))
		return
	}

	// process query params
	query := req.URL.Query()
	log.Printf("Query params: %+v", query)

	ctx := context.Background()
	if timeoutStr := query.Get("timeout"); timeoutStr != "" {
		var cancel context.CancelFunc
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			errorMessage(res, http.StatusBadRequest, fmt.Errorf("timeout duration parsing error: %v", err))
			return
		}
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	var result solveResponse
	defer func() {
		res.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(res)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(result)
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}()

	// solve level
	solver := world.NewSolver(m)
	log.Print("Starting solution...")
	err = solver.Solve(ctx)
	defer func() {
		debug := solver.GetDebug()
		// add additional data
		if query.Get("mapdebug") == "true" {
			result.DeadZones = debug.DeadZones
			result.HorizontalLock = debug.HLock
			result.VerticalLock = debug.VLock
		}
		if query.Get("metrics") == "true" {
			result.Metrics = debug.Metrics
		}
		if query.Get("states") == "true" {
			maxStates := 0
			if countStr := query.Get("max_states"); countStr != "" {
				maxStates, err = strconv.Atoi(countStr)
				if err != nil {
					log.Printf("Error parsing max_states parameter: %v; 1000 will be used as default", err)
					maxStates = 1000
				}
			}
			result.States = solver.GetTree(maxStates)
		}
	}()
	if err != nil {
		log.Printf("Solving error: %v", err)
		result.Error = err.Error()
		return
	}

	// get path
	path, err := solver.GetPath()
	if err != nil {
		log.Printf("Path building error: %v", err)
		result.Error = err.Error()
		return
	}

	result.Solution = path
}

func Run(port int) error {
	http.HandleFunc("/solve", solveHandler)
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}
