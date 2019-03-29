package message

type MessagePageEntry struct {
	Total   int         `json:"total"`
	HasNext bool        `json:"hasNext"`
	Data    interface{} `json:"data"`
}
