package global

import "os"

var (
	MONGO_URL   = "localhost:27017"
	SERVER_PORT = "5000"
)

func init() {
	if mongoUrl := os.Getenv("MONGO_URL"); len(mongoUrl) > 0 {
		MONGO_URL = mongoUrl
	}
	if serverPort := os.Getenv("SERVER_PORT"); len(serverPort) > 0 {
		SERVER_PORT = serverPort
	}
}
