package world

import (
	"encoding/binary"
	"strings"
)

const cellSize = 8

type Bitmap struct {
	cells []uint64
	w, h  int
}

func (b *Bitmap) resizeToFit(p Pos) {
	if b.isInMap(p) {
		return
	}
	cellX := p.X / cellSize
	cellY := p.Y / cellSize
	rw, rh := b.w, b.h
	if cellX+1 > b.w {
		rw = cellX + 1
	}
	if cellY+1 > b.h {
		rh = cellY + 1
	}
	cells := make([]uint64, rw*rh)
	// copy cells
	for i, cell := range b.cells {
		if cell == 0 {
			continue
		}
		ri := rw*(i/b.w) + (i % b.w)
		cells[ri] = cell
	}
	b.w = rw
	b.h = rh
	b.cells = cells
}

func (b Bitmap) cellIndex(p Pos) (int, int) {
	cellX := p.X / cellSize
	cellY := p.Y / cellSize
	cellIndex := cellY*b.w + cellX
	bitX := p.X % cellSize
	bitY := p.Y % cellSize
	bitIndex := bitY*cellSize + bitX
	return cellIndex, bitIndex
}

func (b Bitmap) isInMap(p Pos) bool {
	return p.X >= 0 && p.Y >= 0 && (p.X/cellSize) < b.w && (p.Y/cellSize) < b.h
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

func (b Bitmap) List() []Pos {
	var pos = make([]Pos, 0)
	for i, cell := range b.cells {
		if cell == 0 {
			continue
		}
		cellX := i % b.w
		cellY := i / b.w
		bitX := 0
		bitY := 0
		for cell > 0 {
			if cell&0x1 == 0x1 {
				pos = append(pos, Pos{
					cellX*cellSize + bitX,
					cellY*cellSize + bitY,
				})
			}
			cell >>= 1
			bitX++
			if bitX >= cellSize {
				bitX = 0
				bitY++
			}
		}
	}
	return pos
}

func (b Bitmap) HashBytes() []byte {
	bs := make([]byte, 8*len(b.cells))
	for i, c := range b.cells {
		binary.BigEndian.PutUint64(bs[8*i:], c)
	}
	return bs
}

func (b Bitmap) String() string {
	buf := &strings.Builder{}
	buf.WriteString("Bitmap{")
	for _, pos := range b.List() {
		buf.WriteString(pos.String())
	}
	buf.WriteString("}")
	return buf.String()
}
