package gograph

import (
	"github.com/mtresnik/goutils/pkg/goutils"
	"math/rand"
	"time"
)

type MazeGeneratorRequest struct {
	Rows, Cols          int
	MazeUpdateListeners *[]MazeUpdateListener
}

type MazeUpdateListener interface {
	Update(response MazeGeneratorResponse)
}

func VisitMazeUpdateListeners(listeners []MazeUpdateListener, response MazeGeneratorResponse) {
	for _, listener := range listeners {
		listener.Update(response)
	}
}

type MazeGeneratorResponse struct {
	Maze     Maze
	Visited  map[int64]bool
	Complete bool
}

func NewMazeGeneratorRequest(rows, cols int) MazeGeneratorRequest {
	return MazeGeneratorRequest{
		Rows: rows,
		Cols: cols,
	}
}

type MazeGenerator interface {
	Build(request MazeGeneratorRequest) MazeGeneratorResponse
}

type AldousBroderMazeGenerator struct{}

func (a AldousBroderMazeGenerator) Build(request MazeGeneratorRequest) MazeGeneratorResponse {
	maze := NewMaze(request.Rows, request.Cols)

	mazeUpdateListeners := make([]MazeUpdateListener, 0)
	if request.MazeUpdateListeners != nil {
		mazeUpdateListeners = *request.MazeUpdateListeners
	}

	visitedHashes := make(map[int64]bool)
	visitedCoordinates := make([]MazeCoordinate, 0)
	VisitMazeUpdateListeners(mazeUpdateListeners, MazeGeneratorResponse{maze, visitedHashes, false})

	allCells := maze.Flatten()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(allCells), func(i, j int) {
		allCells[i], allCells[j] = allCells[j], allCells[i]
	})
	randomIndex := random.Intn(len(allCells))
	startCell := &allCells[randomIndex]

	visitedHashes[HashMazeCoordinate(startCell)] = true
	visitedCoordinates = append(visitedCoordinates, startCell)

	VisitMazeUpdateListeners(mazeUpdateListeners, MazeGeneratorResponse{maze, visitedHashes, false})

	numCells := maze.Rows * maze.Cols
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
		neighborConnections = filteredConnections
		if len(neighborConnections) == 0 {
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
					VisitMazeUpdateListeners(mazeUpdateListeners, MazeGeneratorResponse{maze, visitedHashes, true})
					return MazeGeneratorResponse{maze, visitedHashes, true}
				}
				currentCell = found
				visitedHashes[HashMazeCoordinate(currentCell)] = true
				visitedCoordinates = append(visitedCoordinates, currentCell)

				VisitMazeUpdateListeners(mazeUpdateListeners, MazeGeneratorResponse{maze, visitedHashes, false})
				continue
			}
		}
		randomIndex = random.Intn(len(neighborConnections))
		currConnection := neighborConnections[randomIndex]
		maze.SetWall(currConnection.From.GetRow(), currConnection.From.GetCol(), currConnection.Direction, false)
		wall := maze.GetWall(currConnection.From.GetRow(), currConnection.From.GetCol(), currConnection.Direction)
		if wall == nil {
			panic("Wall is nil")
		}
		if *wall {
			panic("Connection is a wall")
		}
		visitedHashes[HashMazeCoordinate(currConnection.To)] = true
		visitedCoordinates = append(visitedCoordinates, currConnection.To)
		currentCell = maze.GetCell(currConnection.To.GetRow(), currConnection.To.GetCol())
		VisitMazeUpdateListeners(mazeUpdateListeners, MazeGeneratorResponse{maze, visitedHashes, false})
	}
	VisitMazeUpdateListeners(mazeUpdateListeners, MazeGeneratorResponse{maze, visitedHashes, true})
	return MazeGeneratorResponse{maze, visitedHashes, true}
}
