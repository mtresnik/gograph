package gomaze

type IMazeToGraphProvider interface {
	Build(maze *Maze)
}

type MazeToGraphProvider struct{}

func (p *MazeToGraphProvider) Build(maze *Maze) {

}
