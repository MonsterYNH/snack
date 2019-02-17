package player

import (
	"errors"
	"snack/node"
)

const (
	DATA_TYPE_NORMAL int = iota
	DATA_TYPE_FOOD
	DATA_TYPE_WALL
	DATA_TYPE_PLAYER

	STATUS_NORMAL
	STATUS_DIE
	STATUS_OPERAT_WRONG

	PLAYER_UP
	PLAYER_DOWN
	PLAYER_RIGHT
	PLAYER_LEFT
)

type Player struct {
	Length int
	Food   int
	Score  int
	Body   *node.LinkNode
	Status int
	Name   string
}

func CreatePlayer(name string, data node.Data) *Player {
	player := &Player{
		Name:   name,
		Length: 1,
		Food:   0,
		Score:  0,
		Status: STATUS_NORMAL,
		Body:   node.CreatNode(),
	}
	player.Body.HeadAddNode(&node.Node{
		Data: data,
	})
	return player
}

func (player *Player) Move(direct, dataType int) (int, error) {
	nodeType := -1
	if err := player.CheckBodyStatus(); err != nil {
		return nodeType, err
	}
	nextNode := player.GetNextNode(direct)
	switch dataType {
	case DATA_TYPE_FOOD:
		player.Body.HeadAddNode(&node.Node{
			Data: nextNode,
		})
		player.Food++
		player.Length++
		player.Score += player.Length * 2
		player.Status = STATUS_NORMAL
		nodeType = DATA_TYPE_FOOD
	case DATA_TYPE_WALL:
		player.Status = STATUS_DIE
		nodeType = DATA_TYPE_WALL
	case DATA_TYPE_PLAYER:
		player.Status = STATUS_DIE
		nodeType = DATA_TYPE_PLAYER
	default:
		player.Body.HeadAddNode(&node.Node{
			Data: nextNode,
		}).TailDeleteNode()
		nodeType = DATA_TYPE_NORMAL
		player.Status = STATUS_NORMAL
	}
	return nodeType, player.CheckBodyStatus()
}

func (player *Player) GetNextNode(direct int) node.Data {
	head := player.Body.Node.Next.Data
	switch direct {
	case PLAYER_UP:
		head.Y++
	case PLAYER_DOWN:
		head.Y--
	case PLAYER_LEFT:
		head.X++
	case PLAYER_RIGHT:
		head.X--
	}
	return head
}

func (player *Player) CheckBodyStatus() error {
	if player.Status == STATUS_DIE {
		return errors.New("player die")
	}
	return nil
}
