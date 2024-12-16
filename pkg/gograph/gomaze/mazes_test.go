package gomaze

import "testing"

func TestMaze_SetWall(t *testing.T) {
	maze := NewMaze(10, 10)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			maze.SetWall(i, j, 0, false)
			wall := maze.GetWall(i, j, 0)
			if wall != nil && *wall {
				t.Error("Wall is a wall")
			}
		}
	}

}
