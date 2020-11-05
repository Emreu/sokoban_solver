package world

import "testing"

func comparePaths(p1, p2 []MoveDirection) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i, dir := range p1 {
		if dir != p2[i] {
			return false
		}
	}
	return true
}

func TestPathfindingSimple(t *testing.T) {
	md := Bitmap{}

	md.SetBit(Pos{0, 0})
	md.SetBit(Pos{1, 0})
	md.SetBit(Pos{2, 0})

	path, err := FindDirections(md, Pos{0, 0}, Pos{2, 0})
	if err != nil {
		t.Fatalf("No simple path found: %v", err)
	}

	expectedPath := []MoveDirection{MoveRight, MoveRight}
	if !comparePaths(path, expectedPath) {
		t.Errorf("Wrong directions, expected: %v got: %v", expectedPath, path)
	}

	path, err = FindDirections(md, Pos{2, 0}, Pos{1, 0})
	if err != nil {
		t.Fatalf("No simple path found: %v", err)
	}

	expectedPath = []MoveDirection{MoveLeft}
	if !comparePaths(path, expectedPath) {
		t.Errorf("Wrong directions, expected: %v got: %v", expectedPath, path)
	}

	md.SetBit(Pos{2, 1})
	md.SetBit(Pos{2, 2})
	md.SetBit(Pos{1, 2})

	path, err = FindDirections(md, Pos{0, 0}, Pos{1, 2})
	if err != nil {
		t.Fatalf("No simple path found: %v", err)
	}

	expectedPath = []MoveDirection{MoveRight, MoveRight, MoveDown, MoveDown, MoveLeft}
	if !comparePaths(path, expectedPath) {
		t.Errorf("Wrong directions, expected: %v got: %v", expectedPath, path)
	}

	path, err = FindDirections(md, Pos{0, 0}, Pos{4, 4})
	if err == nil {
		t.Fatalf("Impossible path found: %v", path)
	}
}
