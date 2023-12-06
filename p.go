package hello

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

type CellName [2]int8
type PieceName [3]int8

type PieceType uint8

const (
	PieceTypeKing PieceType = iota
	PieceTypeQueen
	PieceTypeRook
	PieceTypeBiShop
	PieceTypeKnight
	PieceTypePawn
)

type GameInterface interface {
	IsCellWithinBoard(cellName string) bool
	SetPieceToCell(cellName, pieceName string) bool
	GetAvailableStepsForPiece(cellName string)
}

type stepDelta struct {
	rankDelta int8
	fileDelta int8
}

type PieceKing struct {
	Piece
}

func (p *PieceKing) GetAvailableSteps() []Cell {
	var res []Cell

	// check if piece is dead
	if p.container == nil {
		return res
	}

	predictions := []stepDelta{
		{1, 1}, {-1, -1}, {1, -1}, {-1, 1},
		{1, 0}, {0, 1}, {-1, 0}, {0, -1},
	}
	for _, prediction := range predictions {

	}
}

type PieceQueen struct {
	Piece
}

type PieceRook struct {
	Piece
}

type PieceBiShop struct {
	Piece
}

type PieceKnight struct {
	Piece
}

type PiecePawn struct {
	Piece
}

type Piece struct {
	Type  PieceType
	Color Color
	Index int8

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
	File
	Rank

	piece *Piece // nil means this cell does not contain a piece
}

func (c *Cell) Name() CellName {
	return CellName{int8(c.Rank), int8(c.File)}
}

func (p *Piece) Name() PieceName {
	return PieceName{int8(p.Type), int8(p.Color), p.Index}
}

// A trick
func (c Cell) GetColor() Color {
	sum := uint8(c.File) + uint8(c.Rank)
	return Color(sum & 1)
}

// horizontal row
type Rank int8

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

func (c *Cell) Move(fileDelta, rankDelta int8) Cell {
	return Cell{
		File: File(int8(c.File) + fileDelta),
		Rank: Rank(int8(c.Rank) + rankDelta),
	}
}

// vertical column
type File int8

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

type game struct {
	board  map[CellName]*Cell
	pieces map[PieceName]*Piece
}

func (g *game) IsCellWithinBoard(cellName CellName) bool {
	_, exist := g.board[cellName]
	return exist
}

func (g *game) SetPieceToCell(cellName CellName, pieceName PieceName) bool {
	piece, pieceExist := g.pieces[pieceName]
	cell, cellExist := g.board[cellName]

	if cellExist && pieceExist && cell != nil && piece != nil {
		piece.container = cell

		if cell.piece != nil {
			delete(g.pieces, cell.piece.Name())
		}
		cell.piece = piece
		return true
	}
	return false
}

func (g *game) GetAvailableStepsForPiece() []Cell {

}

func createGame() GameInterface {

}
