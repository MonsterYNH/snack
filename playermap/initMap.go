package playermap

var playerMap *PlayerMap

func init() {
	playerMap = CreatePlayerMap(100, 100)
}

func GetMap() *PlayerMap {
	return playerMap
}

func GetClient(name string) (*ClientInfo, error) {
	return GetClientInfo(name)
}

func Play(operat int, playerName string) error {
	return playerMap.Play(Operat{
		Operator:   operat,
		PlayerName: playerName,
	})
}
