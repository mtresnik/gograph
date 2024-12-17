package gograph

import "testing"

func TestAldousBroderMazeGenerator_Build(t *testing.T) {
	AldousBroderMazeGenerator{}.Build(NewMazeGeneratorRequest(10, 10))
}
