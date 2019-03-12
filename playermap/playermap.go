package playermap

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"snack/db"
	"snack/node"
	"snack/player"
	"time"

	"github.com/garyburd/redigo/redis"
)

type PlayerMap struct {
	Name       string                    `json:"name"`
	Players    map[string]*player.Player `json:"players"`
	PlayMap    map[int][]*node.Data      `json:"play_map"`
	PlayerChan chan *player.Player       `json:"-"`
	OperatChan chan interface{}          `json:"-"`
	Width      int                       `json:"width"`
	Length     int                       `json:"length"`
}

type Operat struct {
	Operator   int
	PlayerName string
}

type DataEntry struct {
	Node node.Data
	Type int
}

type ClientInfo struct {
	Mine   *player.Player       `json:"mine"`
	Others []*player.Player     `json:"others"`
	Map    map[int][]*node.Data `json:"map"`
}

func CreatePlayerMap(width, length int) *PlayerMap {
	playerMap := &PlayerMap{
		Players:    make(map[string]*player.Player),
		PlayMap:    make(map[int][]*node.Data),
		PlayerChan: make(chan *player.Player),
		OperatChan: make(chan interface{}),
		Width:      width,
		Length:     length,
	}
	playerMap.addFoods(10)
	go func(playerMap *PlayerMap) {
		go func() {
			if err := recover(); err != nil {
				fmt.Println(err, "=====")
			}
		}()
		conn := db.RedisPool.Get()
		defer conn.Close()

		_, err := conn.Do("SET", "play_map", "")
		fmt.Println(err)
		for {
			select {
			case data := <-playerMap.OperatChan: // 所有操作PlayMap属性都必须在此
				if dataEntry, exist := data.(DataEntry); exist {
					playerMap.updateNode(dataEntry.Type, &dataEntry.Node)
				} else if operat, exist := data.(Operat); exist {
					if err := playerMap.move(operat.PlayerName, operat.Operator); err != nil {
						playerMap.bodyConvertFood(operat.PlayerName)
					}
				}
			case data := <-playerMap.PlayerChan: // 所有用户操作都必须在此
				playerMap.Players[data.Name] = data
			}
			// 将地图数据写入Redis
			bytes, _ := json.Marshal(playerMap)

			if _, err := conn.Do("SET", "play_map", bytes); err != nil {
				fmt.Println("Redis Error:", err.Error())
			}
		}
	}(playerMap)
	return playerMap
}

func (playerMap *PlayerMap) AddPlayer(user *player.Player) {
	playerMap.PlayerChan <- user
	playerMap.OperatChan <- DataEntry{
		Node: user.Body.Node.Next.Data,
		Type: player.DATA_TYPE_PLAYER,
	}
}

func (playerMap *PlayerMap) Play(operat Operat) error {
	if err := playerMap.Players[operat.PlayerName].CheckBodyStatus(); err != nil {
		return err
	}
	playerMap.OperatChan <- operat
	return nil
}

// playMap
func GetClientInfo(name string) (*ClientInfo, error) {
	conn := db.RedisPool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", "play_map"))
	if err != nil {
		return nil, err
	}

	playerMap := PlayerMap{}
	if err := json.Unmarshal(data, &playerMap); err != nil {
		return nil, err
	}
	mine := playerMap.Players[name]
	others := make([]*player.Player, len(playerMap.Players)-1)
	for _, user := range playerMap.Players {
		if name != user.Name {
			others = append(others, user)
		}
	}
	clientInfo := &ClientInfo{
		Mine:   mine,
		Others: others,
		//Foods:  len(playerMap.PlayMap[player.DATA_TYPE_FOOD]),
		Map: playerMap.PlayMap,
	}
	return clientInfo, nil
}

// playMap read
func GetMapInfo() {
	conn := db.RedisPool.Get()
	defer conn.Close()

	data, _ := redis.Bytes(conn.Do("GET", "play_map"))
	playerMap := PlayerMap{}
	if err := json.Unmarshal(data, &playerMap); err != nil {
		return
	}
	fmt.Println("Foods: ", len(playerMap.PlayMap[player.DATA_TYPE_FOOD]), " User: ", len(playerMap.Players))
	bytes, _ := json.Marshal(playerMap.PlayMap[player.DATA_TYPE_PLAYER])
	fmt.Println("Player: ", string(bytes))
}

func (playerMap *PlayerMap) ReLive(name string) {}

// playMap
func (playerMap *PlayerMap) updateNode(nodeType int, data *node.Data) {
	switch nodeType {
	case player.DATA_TYPE_FOOD:
		if playerMap.PlayMap[player.DATA_TYPE_FOOD] == nil {
			playerMap.PlayMap[player.DATA_TYPE_FOOD] = make([]*node.Data, 0)
		}
		playerMap.PlayMap[player.DATA_TYPE_FOOD] = append(playerMap.PlayMap[player.DATA_TYPE_FOOD], data)
	case player.DATA_TYPE_WALL:
		if playerMap.PlayMap[player.DATA_TYPE_WALL] == nil {
			playerMap.PlayMap[player.DATA_TYPE_WALL] = make([]*node.Data, 0)
		}
		playerMap.PlayMap[player.DATA_TYPE_WALL] = append(playerMap.PlayMap[player.DATA_TYPE_WALL], data)
	case player.DATA_TYPE_PLAYER:
		if playerMap.PlayMap[player.DATA_TYPE_PLAYER] == nil {
			playerMap.PlayMap[player.DATA_TYPE_PLAYER] = make([]*node.Data, 0)
		}
		playerMap.PlayMap[player.DATA_TYPE_PLAYER] = append(playerMap.PlayMap[player.DATA_TYPE_PLAYER], data)
	}
}

// playMap
func (playerMap *PlayerMap) move(name string, operator int) error {
	user := playerMap.Players[name]
	if err := user.CheckBodyStatus(); err != nil {
		return err
	}
	// 计算落点
	nextNode := user.GetNextNode(operator)
	// 判断越界
	if playerMap.checkBorder(nextNode) {
		user.Status = player.STATUS_DIE
	}
	// 落点的类型
	nodeType, _ := playerMap.getNodeType(nextNode)
	if _, err := user.Move(operator, nodeType); err != nil {
		user.Status = player.STATUS_DIE
	} else {
		playerMap.updatePlayerInMap()
		if isEate := playerMap.eateFood(nextNode); isEate {
			playerMap.addFoods(1)
		}
	}
	return user.CheckBodyStatus()
}

func (playerMap *PlayerMap) createRandomNode() node.Data {
	rand.Seed(time.Now().UnixNano())
	return node.Data{
		X: rand.Intn(playerMap.Length + 1),
		Y: rand.Intn(playerMap.Width + 1),
	}
}

func (playerMap *PlayerMap) addFoods(num int) {
	for i := 0; i < num; i++ {
		data := playerMap.createRandomNode()
		if nodeType, _ := playerMap.getNodeType(data); nodeType < 0 {
			playerMap.PlayMap[player.DATA_TYPE_FOOD] = append(playerMap.PlayMap[player.DATA_TYPE_FOOD], &data)
		} else {
			i--
		}
	}
}

func (playerMap *PlayerMap) updatePlayerInMap() {
	playerNodes := make([]*node.Data, 0)
	for _, user := range playerMap.Players {
		nodeEntry := user.Body.Node.Next
		for {
			playerNodes = append(playerNodes, &nodeEntry.Data)
			if nodeEntry.Next == nil {
				break
			}
			nodeEntry = nodeEntry.Next
		}
	}
	playerMap.PlayMap[player.DATA_TYPE_PLAYER] = playerNodes
}

func (playerMap *PlayerMap) checkBorder(data node.Data) bool {
	if data.X < 0 || data.X > playerMap.Length {
		return true
	}
	if data.Y < 0 || data.Y > playerMap.Width {
		return true
	}
	return false
}

// playMap
func (playerMap *PlayerMap) getNodeType(data node.Data) (int, int) {
	for nodeType, nodes := range playerMap.PlayMap {
		for index, nodeEntry := range nodes {
			if nodeEntry.X == data.X && nodeEntry.Y == data.Y {
				return nodeType, index
			}
		}
	}
	return -1, -1
}

func (playerMap *PlayerMap) eateFood(data node.Data) bool {
	nodeType, index := playerMap.getNodeType(data)
	if nodeType == player.DATA_TYPE_FOOD && index >= 0 {
		foods := playerMap.PlayMap[player.DATA_TYPE_FOOD]
		playerMap.PlayMap[player.DATA_TYPE_FOOD] = append(foods[:index], foods[index+1:]...)
		return true
	}
	return false
}

func (playerMap *PlayerMap) bodyConvertFood(name string) {
	data := playerMap.Players[name].Body.Node.Next
	for {
		userBody := playerMap.PlayMap[player.DATA_TYPE_PLAYER]
		for index, nodeEntry := range userBody {
			if nodeEntry.X == data.Data.X && nodeEntry.Y == data.Data.Y {
				userBody = append(userBody[:index], userBody[index+1:]...)
			}
		}
		playerMap.PlayMap[player.DATA_TYPE_PLAYER] = userBody
		playerMap.PlayMap[player.DATA_TYPE_FOOD] = append(playerMap.PlayMap[player.DATA_TYPE_FOOD], &data.Data)
		if data.Next == nil {
			break
		}
		data = data.Next
	}

}

// playMap read
func PrintMap() {
	conn := db.RedisPool.Get()
	defer conn.Close()

	data, _ := redis.Bytes(conn.Do("GET", "play_map"))
	playerMap := PlayerMap{}
	if err := json.Unmarshal(data, &playerMap); err != nil {
		return
	}
	for y := 0; y < playerMap.Width; y++ {
		fmt.Print("|")
		for x := 0; x < playerMap.Length; x++ {
			nodeType, _ := playerMap.getNodeType(node.Data{
				X: x,
				Y: y,
			})
			if nodeType == player.DATA_TYPE_FOOD {
				fmt.Print("*")
			} else if nodeType == player.DATA_TYPE_WALL {
				fmt.Print("#")
			} else if nodeType == player.DATA_TYPE_PLAYER {
				fmt.Print("@")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("|")
	}
}
