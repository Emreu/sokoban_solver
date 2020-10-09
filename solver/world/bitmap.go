package world

const cellSize = 8

type Bitmap struct {
	cells []uint64
	w, h  int
}

func (b *Bitmap) resizeToFit(p Pos) {
	if b.isInMap(p) {
		return
	}
	// TODO: implement
}

func (b Bitmap) cellIndex(p Pos) (int, int) {
	return 0, 0
}

func (b Bitmap) isInMap(p Pos) bool {
	return p.X > 0 && p.Y > 0 && (p.X/cellSize) < b.w && (p.Y/cellSize) < b.h
}

func (b Bitmap) CheckBit(p Pos) bool {
	if !b.isInMap(p) {
		return false
	}
	c, bit := b.cellIndex(p)
	cell := b.cells[c]
	return (cell>>bit)&0x1 == 0x1
}

func (b *Bitmap) SetBit(p Pos) {
	b.resizeToFit(p)
	c, bit := b.cellIndex(p)
	b.cells[c] |= 0x1 << bit
}
