package world

import "testing"

func TestBitmapBasic(t *testing.T) {
	b := Bitmap{}

	if b.CheckBit(Pos{0, 0}) {
		t.Error("Empty bitmap bit @(0,0) set")
	}

	b.SetBit(Pos{0, 0})

	if !b.CheckBit(Pos{0, 0}) {
		t.Error("Bitmap bit @(0,0) not checked after setting")
	}

	b.SetBit(Pos{9, 9})

	if !b.CheckBit(Pos{9, 9}) {
		t.Error("Bitmap bit @(9,9) not checked after setting")
	}
}

func TestBitmapListPositions(t *testing.T) {
	b := Bitmap{}

	pos := b.List()
	if len(pos) > 0 {
		t.Error("Empty bitmap listed non-zero positions")
	}

	b.SetBit(Pos{0, 0})
	b.SetBit(Pos{8, 8})

	pos = b.List()
	if len(pos) != 2 {
		t.Fatalf("Bitmap with 2 bits listed %d positions", len(pos))
	}

	if pos[0].X != 0 || pos[0].Y != 0 {
		t.Fatalf("1st position is not @(0,0) but %s", pos[0])
	}

	if pos[1].X != 8 || pos[1].Y != 8 {
		t.Fatalf("2nd position is not @(8,8) but %s", pos[1])
	}

	b.SetBit(Pos{3, 5})
	b.SetBit(Pos{7, 11})

	pos = b.List()
	if len(pos) != 4 {
		t.Fatalf("Bitmap with 4 bits listed %d positions", len(pos))
	}

	if pos[1].X != 3 || pos[1].Y != 5 {
		t.Fatalf("2nd position is not @(3,5) but %s", pos[1])
	}

	if pos[2].X != 7 || pos[2].Y != 11 {
		t.Fatalf("3rd position is not @(7,11) but %s", pos[2])
	}
}
