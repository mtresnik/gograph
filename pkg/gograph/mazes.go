package gograph

import (
	"hash/fnv"
	"math"
)

const (
	DIRECTION_RIGHT = 0
	DIRECTION_UP    = 1
	DIRECTION_LEFT  = 2
	DIRECTION_DOWN  = 3
)

type MazeCoordinate interface {
	GetRow() int
	GetCol() int
}

func HashMazeCoordinate(c MazeCoordinate) int64 {
	values := []float64{float64(c.GetRow()), float64(c.GetCol())}
	hasher := fnv.New64a()

	for _, value := range values {
		bits := math.Float64bits(value)
		buf := make([]byte, 8)
		for i := 0; i < 8; i++ {
			buf[i] = byte(bits >> (i * 8))
		}
		_, _ = hasher.Write(buf)
	}

	return int64(hasher.Sum64())
}

func EqualsMazeCoordinate(a, b MazeCoordinate) bool {
	return HashMazeCoordinate(a) == HashMazeCoordinate(b)
}

type MazeCell struct {
	Row        int
	Col        int
	LeftIsWall bool
	UpIsWall   bool
}

func (m MazeCell) GetRow() int {
	return m.Row
}

func (m MazeCell) GetCol() int {
	return m.Col
}

type MazeConnection struct {
	From      MazeCoordinate
	To        MazeCoordinate
	Direction int
	IsWall    bool
}

func (m MazeConnection) Equals(other MazeConnection) bool {
	return EqualsMazeCoordinate(m.From, other.From) && EqualsMazeCoordinate(m.To, other.To) ||
		EqualsMazeCoordinate(m.From, other.To) && EqualsMazeCoordinate(m.To, other.From)
}

func (m MazeCell) Equals(other MazeCell) bool {
	return m.Row == other.Row && m.Col == other.Col
}

func (m MazeCell) Hash() int64 {
	return HashMazeCoordinate(m)
}

type Maze struct {
	Rows  int
	Cols  int
	Cells [][]MazeCell
}

func NewMaze(rows, cols int) Maze {
	cells := make([][]MazeCell, rows)
	for i := 0; i < rows; i++ {
		cells[i] = make([]MazeCell, cols)
		for j := 0; j < cols; j++ {
			cells[i][j] = MazeCell{
				Row:        i,
				Col:        j,
				LeftIsWall: true,
				UpIsWall:   true,
			}
		}
	}
	return Maze{
		Rows:  rows,
		Cols:  cols,
		Cells: cells,
	}
}

func (m Maze) Flatten() []MazeCell {
	retArray := make([]MazeCell, 0)
	for _, row := range m.Cells {
		retArray = append(retArray, row...)
	}
	return retArray
}

func (m Maze) Contains(row, col int) bool {
	return row >= 0 && row < m.Rows && col >= 0 && col < m.Cols
}

func (m Maze) GetCell(row, col int) *MazeCell {
	return &m.Cells[row][col]
}

func (m Maze) GetIfValid(row, col int) *MazeCell {
	if m.Contains(row, col) {
		return &m.Cells[row][col]
	}
	return nil
}

func (m Maze) GetRight(row, col int) *MazeCell {
	if m.Contains(row, col+1) {
		return &m.Cells[row][col+1]
	}
	return nil
}

func (m Maze) GetUp(row, col int) *MazeCell {
	if m.Contains(row-1, col) {
		return &m.Cells[row-1][col]
	}
	return nil
}

func (m Maze) GetLeft(row, col int) *MazeCell {
	if m.Contains(row, col-1) {
		return &m.Cells[row][col-1]
	}
	return nil
}

func (m Maze) GetDown(row, col int) *MazeCell {
	if m.Contains(row+1, col) {
		return &m.Cells[row+1][col]
	}
	return nil
}

func (m Maze) GetNeighbor(row, col int, direction int) *MazeCell {
	switch direction {
	case DIRECTION_RIGHT:
		return m.GetRight(row, col)
	case DIRECTION_UP:
		return m.GetUp(row, col)
	case DIRECTION_LEFT:
		return m.GetLeft(row, col)
	case DIRECTION_DOWN:
		return m.GetDown(row, col)
	}
	return nil
}

func (m Maze) GetNeighbors(row, col int) []MazeCell {
	retArray := make([]MazeCell, 0)
	right := m.GetRight(row, col)
	up := m.GetUp(row, col)
	left := m.GetLeft(row, col)
	down := m.GetDown(row, col)
	if right != nil {
		retArray = append(retArray, *right)
	}
	if up != nil {
		retArray = append(retArray, *up)
	}
	if left != nil {
		retArray = append(retArray, *left)
	}
	if down != nil {
		retArray = append(retArray, *down)
	}
	return retArray
}

func (m Maze) GetConnections(row, col int) []MazeConnection {
	retArray := make([]MazeConnection, 0)
	if !m.Contains(row, col) {
		return retArray
	}
	directions := []int{DIRECTION_RIGHT, DIRECTION_UP, DIRECTION_LEFT, DIRECTION_DOWN}
	for _, direction := range directions {
		connection := m.GetConnection(row, col, direction)
		if connection != nil {
			retArray = append(retArray, *connection)
		}
	}
	return retArray
}

func (m Maze) GetConnection(row, col int, direction int) *MazeConnection {
	from := m.GetIfValid(row, col)
	if from == nil {
		return nil
	}
	to := m.GetNeighbor(row, col, direction)
	if to == nil {
		return nil
	}
	wall := m.GetWall(row, col, direction)
	if wall == nil {
		return nil
	}
	return &MazeConnection{
		From:      *from,
		To:        *to,
		Direction: direction,
		IsWall:    *wall,
	}
}

func (m Maze) GetWall(row, col int, direction int) *bool {
	switch direction {
	case DIRECTION_RIGHT:
		return m.GetRightWall(row, col)
	case DIRECTION_UP:
		return m.GetUpWall(row, col)
	case DIRECTION_LEFT:
		return m.GetLeftWall(row, col)
	case DIRECTION_DOWN:
		return m.GetDownWall(row, col)
	default:
		panic("Invalid direction")
	}
}

func (m Maze) GetRightWall(row, col int) *bool {
	if !m.Contains(row, col+1) {
		return nil
	}
	rightCell := m.GetRight(row, col)
	if rightCell == nil {
		return nil
	}
	return &rightCell.LeftIsWall
}

func (m Maze) GetUpWall(row, col int) *bool {
	if !m.Contains(row, col) {
		return nil
	}
	cell := m.GetCell(row, col)
	return &cell.UpIsWall
}

func (m Maze) GetLeftWall(row, col int) *bool {
	if !m.Contains(row, col) {
		return nil
	}
	cell := m.GetCell(row, col)
	return &cell.LeftIsWall
}

func (m Maze) GetDownWall(row, col int) *bool {
	if !m.Contains(row+1, col) {
		return nil
	}
	cell := m.GetDown(row, col)
	if cell == nil {
		return nil
	}
	return &cell.UpIsWall
}

func (m Maze) SetWall(row, col int, direction int, border bool) {
	if !m.Contains(row, col) {
		return
	}
	switch direction {
	case DIRECTION_RIGHT:
		m.SetRightWall(row, col, border)
		return
	case DIRECTION_UP:
		m.SetUpWall(row, col, border)
		return
	case DIRECTION_LEFT:
		m.SetLeftWall(row, col, border)
		return
	case DIRECTION_DOWN:
		m.SetDownWall(row, col, border)
		return
	}
}

func (m Maze) SetRightWall(row, col int, border bool) {
	if !m.Contains(row, col+1) {
		return
	}
	cell := m.GetCell(row, col+1)
	cell.LeftIsWall = border
}

func (m Maze) SetUpWall(row, col int, border bool) {
	if !m.Contains(row, col) {
		return
	}
	cell := m.GetCell(row, col)
	cell.UpIsWall = border
}

func (m Maze) SetLeftWall(row, col int, border bool) {
	if !m.Contains(row, col) {
		return
	}
	cell := m.GetCell(row, col)
	cell.LeftIsWall = border
}

func (m Maze) SetDownWall(row, col int, border bool) {
	if !m.Contains(row+1, col) {
		return
	}
	cell := m.GetCell(row+1, col)
	cell.UpIsWall = border
}
