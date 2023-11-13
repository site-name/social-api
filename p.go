package main

/*
  1 2 3 4 5 6 7 8
a o x o x o x o x
b x o x o x o x o
c o x o x o x o x
d x o x o x o x o
e o x o x o x o x
f x o x o x o x o
g o x o x o x o x
h x o x o x o x o
*/

type PieceType uint8

const (
	PieceTypeKing PieceType = iota
	PieceTypeQueen
	PieceTypeRook
	PieceTypeBiShop
	PieceTypeKnight
	PieceTypePawn
)

// type Side uint8

// const (
// 	SideGood Side = iota
// 	SideBad
// )

type Piece struct {
	Type  PieceType
	Color Color

	container *Cell // nil means this piece is dead
}

func (p Piece) IsDead() bool {
	return p.container == nil
}

type Color uint8

const (
	Black = iota
	White
)

type Cell struct {
	File File
	Rank Rank

	piece *Piece // nil means this cell does not contain a piece
}

func (c *Cell) SetPiece(p *Piece) {
	c.piece = p
	if p != nil {
		p.container = c
	}
}

func (c Cell) GetPiece() *Piece {
	return c.piece
}

// A trick
func (c Cell) GetColor() Color {
	sum := uint8(c.File) + uint8(c.Rank)
	return Color(sum & 1)
}

// horizontal row
type Rank uint8

const (
	Rank1 Rank = 1 + iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

// vertical column
type File byte

const (
	FileA File = 'a'
	FileB File = 'b'
	FileC File = 'c'
	FileD File = 'd'
	FileE File = 'e'
	FileF File = 'f'
	FileG File = 'g'
	FileH File = 'h'
)

// var Coordinate struct {
// 	Rank Rank
// 	File File
// }

func hi() {
	cell_A1.GetPiece().IsDead()
}

var (
	cell_A1 = Cell{FileA, Rank1}
	cell_A2 = Cell{FileA, Rank2}
	cell_A3 = Cell{FileA, Rank3}
	cell_A4 = Cell{FileA, Rank4}
	cell_A5 = Cell{FileA, Rank5}
	cell_A6 = Cell{FileA, Rank6}
	cell_A7 = Cell{FileA, Rank7}
	cell_A8 = Cell{FileA, Rank8}

	cell_B1 = Cell{FileB, Rank1}
	cell_B2 = Cell{FileB, Rank2}
	cell_B3 = Cell{FileB, Rank3}
	cell_B4 = Cell{FileB, Rank4}
	cell_B5 = Cell{FileB, Rank5}
	cell_B6 = Cell{FileB, Rank6}
	cell_B7 = Cell{FileB, Rank7}
	cell_B8 = Cell{FileB, Rank8}

	cell_C1 = Cell{FileC, Rank1}
	cell_C2 = Cell{FileC, Rank2}
	cell_C3 = Cell{FileC, Rank3}
	cell_C4 = Cell{FileC, Rank4}
	cell_C5 = Cell{FileC, Rank5}
	cell_C6 = Cell{FileC, Rank6}
	cell_C7 = Cell{FileC, Rank7}
	cell_C8 = Cell{FileC, Rank8}

	cell_D1 = Cell{FileD, Rank1}
	cell_D2 = Cell{FileD, Rank2}
	cell_D3 = Cell{FileD, Rank3}
	cell_D4 = Cell{FileD, Rank4}
	cell_D5 = Cell{FileD, Rank5}
	cell_D6 = Cell{FileD, Rank6}
	cell_D7 = Cell{FileD, Rank7}
	cell_D8 = Cell{FileD, Rank8}

	cell_E1 = Cell{FileE, Rank1}
	cell_E2 = Cell{FileE, Rank2}
	cell_E3 = Cell{FileE, Rank3}
	cell_E4 = Cell{FileE, Rank4}
	cell_E5 = Cell{FileE, Rank5}
	cell_E6 = Cell{FileE, Rank6}
	cell_E7 = Cell{FileE, Rank7}
	cell_E8 = Cell{FileE, Rank8}

	cell_F1 = Cell{FileF, Rank1}
	cell_F2 = Cell{FileF, Rank2}
	cell_F3 = Cell{FileF, Rank3}
	cell_F4 = Cell{FileF, Rank4}
	cell_F5 = Cell{FileF, Rank5}
	cell_F6 = Cell{FileF, Rank6}
	cell_F7 = Cell{FileF, Rank7}
	cell_F8 = Cell{FileF, Rank8}

	cell_G1 = Cell{FileG, Rank1}
	cell_G2 = Cell{FileG, Rank2}
	cell_G3 = Cell{FileG, Rank3}
	cell_G4 = Cell{FileG, Rank4}
	cell_G5 = Cell{FileG, Rank5}
	cell_G6 = Cell{FileG, Rank6}
	cell_G7 = Cell{FileG, Rank7}
	cell_G8 = Cell{FileG, Rank8}

	cell_H1 = Cell{FileH, Rank1}
	cell_H2 = Cell{FileH, Rank2}
	cell_H3 = Cell{FileH, Rank3}
	cell_H4 = Cell{FileH, Rank4}
	cell_H5 = Cell{FileH, Rank5}
	cell_H6 = Cell{FileH, Rank6}
	cell_H7 = Cell{FileH, Rank7}
	cell_H8 = Cell{FileH, Rank8}
)

// var initialSetup map[]
