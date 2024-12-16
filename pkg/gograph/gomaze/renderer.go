package gomaze

import (
	"github.com/mtresnik/goutils/pkg/goutils"
	"image"
	"image/color"
	"image/gif"
)

type IMazeRenderer interface {
	RenderFrame(m MazeGeneratorResponse) *image.Paletted
}

type ILiveMazeRenderer interface {
	RenderFrames() *gif.GIF
}

type LiveMazeRenderer struct {
	CellSize     int
	VisitedColor *color.RGBA
	Request      MazeGeneratorRequest
	Generator    MazeGenerator
	Frames       []*image.Paletted
}

func (m *LiveMazeRenderer) Build() {
	m.Request.MazeUpdateListeners = &[]MazeUpdateListener{m}
	m.Generator.Build(m.Request)
}

func (m *LiveMazeRenderer) OnMazeUpdated(response MazeGeneratorResponse) {
	internalMazeRenderer := MazeRenderer{CellSize: m.CellSize, VisitedColor: m.VisitedColor, Request: m.Request, Generator: m.Generator}
	frame := internalMazeRenderer.RenderFrame(response)
	m.Frames = append(m.Frames, frame)
}

func (m *LiveMazeRenderer) RenderFrames() *gif.GIF {
	m.Build()
	images := m.Frames
	delays := make([]int, len(images))
	for i := 0; i < len(images); i++ {
		delays[i] = 1
	}
	g := &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: -1,
	}
	return g
}

type MazeRenderer struct {
	CellSize     int
	VisitedColor *color.RGBA
	Request      MazeGeneratorRequest
	Generator    MazeGenerator
}

func (m MazeRenderer) Build() *image.Paletted {
	mazeResponse := m.Generator.Build(m.Request)
	return m.RenderFrame(mazeResponse)
}

func (m MazeRenderer) drawMazeCell(mazeResponse MazeGeneratorResponse, imageColStart, imageColEnd, imageRowStart, imageRowEnd int, image *image.Paletted, mazeRow, mazeCol int) {
	wallSize := m.CellSize / 8
	cell := mazeResponse.Maze.GetCell(mazeRow, mazeCol)
	visitedColor := m.VisitedColor
	if visitedColor == nil {
		white := color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255} // White
		visitedColor = &white
	}
	maze := mazeResponse.Maze
	// Set background for cell if visited
	if goutils.SetContains(mazeResponse.Visited, cell.Hash()) {
		for imageCol := imageColStart; imageCol < imageColEnd; imageCol++ {
			for imageRow := imageRowStart; imageRow < imageRowEnd; imageRow++ {
				image.Set(imageCol, imageRow, *visitedColor)
			}
		}
	}
	leftWall := maze.GetLeftWall(mazeRow, mazeCol)
	rightWall := maze.GetRightWall(mazeRow, mazeCol)
	upWall := maze.GetUpWall(mazeRow, mazeCol)
	downWall := maze.GetDownWall(mazeRow, mazeCol)
	if leftWall != nil && *leftWall {
		for imageRow := imageRowStart; imageRow < imageRowEnd; imageRow++ {
			for wallIndex := 0; wallIndex < wallSize; wallIndex++ {
				image.SetColorIndex(imageColStart+wallIndex, imageRow, 1)
			}
		}
	}
	if rightWall != nil && *rightWall {
		for imageRow := imageRowStart; imageRow < imageRowEnd; imageRow++ {
			for wallIndex := 0; wallIndex < wallSize; wallIndex++ {
				image.SetColorIndex(imageColEnd-wallIndex-1, imageRow, 1)
			}
		}
	}
	if upWall != nil && *upWall {
		for imageCol := imageColStart; imageCol < imageColEnd; imageCol++ {
			for wallIndex := 0; wallIndex < wallSize; wallIndex++ {
				image.SetColorIndex(imageCol, imageRowStart-wallIndex-1, 1)
			}
		}
	}
	if downWall != nil && *downWall {
		for imageCol := imageColStart; imageCol < imageColEnd; imageCol++ {
			for wallIndex := 0; wallIndex < wallSize; wallIndex++ {
				image.SetColorIndex(imageCol, imageRowEnd+wallIndex, 1)
			}
		}
	}
}

func (m MazeRenderer) RenderFrame(response MazeGeneratorResponse) *image.Paletted {
	height := m.CellSize * response.Maze.Rows
	width := m.CellSize * response.Maze.Cols
	visitedColor := m.VisitedColor
	if visitedColor == nil {
		visitedColor = &color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}

	palette := color.Palette{
		color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White
		color.RGBA{0, 0, 0, 255},                   //Black
		color.RGBA{255, 0, 0, 255},                 // Red
		color.RGBA{
			R: 100,
			G: 100,
			B: 255,
			A: 255,
		}, // Blue
		*visitedColor,
	}
	retImage := image.NewPaletted(image.Rect(0, 0, width, height), palette)
	// Clear canvas
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			retImage.SetColorIndex(col, row, 0)
		}
	}
	maze := response.Maze
	for mazeRow := 0; mazeRow < maze.Rows; mazeRow++ {
		imageRowStart := mazeRow * m.CellSize
		imageRowEnd := (mazeRow + 1) * m.CellSize
		for mazeCol := 0; mazeCol < maze.Cols; mazeCol++ {
			imageColStart := mazeCol * m.CellSize
			imageColEnd := (mazeCol + 1) * m.CellSize
			m.drawMazeCell(response, imageColStart, imageColEnd, imageRowStart, imageRowEnd, retImage, mazeRow, mazeCol)
		}
	}
	return retImage
}
