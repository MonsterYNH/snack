package player

import (
	"snack/node"
	"testing"
)

func TestSnack(t *testing.T) {
	player := CreatePlayer("test", node.Data{
		X: 5,
		Y: 5,
	})
	player.Move(PLAYER_UP, DATA_TYPE_FOOD)
	player.Body.PrintLinkNode()
}
