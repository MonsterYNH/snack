package node

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	snack := CreatNode()
	snack.HeadAddNode(&Node{
		Data: Data{
			X: 1,
			Y: 1,
		},
	}).HeadAddNode(&Node{
		Data: Data{
			X: 2,
			Y: 2,
		},
	}).HeadAddNode(&Node{
		Data: Data{
			X: 3,
			Y: 3,
		},
	}).TailDeleteNode().TailDeleteNode().TailDeleteNode().TailDeleteNode()
	bytes, _ := json.Marshal(snack)
	fmt.Println(string(bytes))
}
