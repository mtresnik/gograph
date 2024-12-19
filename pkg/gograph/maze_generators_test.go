package gograph

import "testing"

func TestAldousBroderMazeGenerator_Build(t *testing.T) {
	AldousBroderMazeGenerator(NewMazeGeneratorRequest(10, 10))
}
