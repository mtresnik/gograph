package gograph

import (
	"github.com/mtresnik/goutils/pkg/goutils"
	"math/rand"
	"time"
)

type MazeGenerator interface {
	Build(rows, cols int) Maze
}

type AldousBroderMazeGenerator struct{}

func (a AldousBroderMazeGenerator) Build(rows, cols int) Maze {
	maze := NewMaze(rows, cols)

	visitedHashes := make(map[int64]bool)
	visitedCoordinates := make([]MazeCoordinate, 0)

	allCells := maze.Flatten()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(allCells), func(i, j int) {
		allCells[i], allCells[j] = allCells[j], allCells[i]
	})
	randomIndex := random.Intn(len(allCells))
	startCell := allCells[randomIndex]

	visitedHashes[HashMazeCoordinate(startCell)] = true
	visitedCoordinates = append(visitedCoordinates, startCell)

	numCells := len(allCells)
	currentCell := startCell
	for len(visitedHashes) != numCells {
		neighborConnections := maze.GetConnections(currentCell.Row, currentCell.Col)
		filteredConnections := make([]MazeConnection, 0)
		for _, connection := range neighborConnections {
			hashed := HashMazeCoordinate(connection.To)
			if !goutils.SetContains(visitedHashes, hashed) {
				filteredConnections = append(filteredConnections, connection)
			}
		}
		if len(filteredConnections) == 0 {
			// Attempt to iterate over visited coordinates to find an open path somewhere
			neighborConnections = make([]MazeConnection, 0)
			for _, mazeCoordinate := range visitedCoordinates {
				currConnections := maze.GetConnections(mazeCoordinate.GetRow(), mazeCoordinate.GetCol())
				for _, connection := range currConnections {
					hashedConnectionTo := HashMazeCoordinate(connection.To)
					if !goutils.SetContains(visitedHashes, hashedConnectionTo) {
						neighborConnections = append(neighborConnections, connection)
					}
				}
			}
			if len(neighborConnections) == 0 {
				// Find the first non visited cell
				var found *MazeCell = nil
				for _, cell := range allCells {
					if !goutils.SetContains(visitedHashes, HashMazeCoordinate(cell)) {
						found = &cell
						break
					}
				}
				if found == nil {
					return maze
				}
				currentCell = *found
				visitedHashes[HashMazeCoordinate(currentCell)] = true
				visitedCoordinates = append(visitedCoordinates, currentCell)
				continue
			}
		}
		randomIndex = random.Intn(len(neighborConnections))
		currConnection := neighborConnections[randomIndex]
		maze.SetWall(currConnection.From.GetRow(), currConnection.From.GetCol(), currConnection.Direction, MazeBorder{IsWall: false})
		visitedHashes[HashMazeCoordinate(currConnection.To)] = true
		visitedCoordinates = append(visitedCoordinates, currConnection.To)
		currentCell = maze.GetCell(currConnection.To.GetRow(), currConnection.To.GetCol())
	}
	return maze
}
