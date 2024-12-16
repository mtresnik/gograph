package gograph

import "testing"

func TestAldousBroderMazeGenerator_Build(t *testing.T) {
	AldousBroderMazeGenerator{}.Build(10, 10)
}
