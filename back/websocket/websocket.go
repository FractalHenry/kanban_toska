package websocket

/* import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type SpaceConnection struct {
	SpaceID uint
	Conn    *websocket.Conn
}

type BoardConnection struct {
	SpaceID uint
	BoardID uint
	Conn    *websocket.Conn
}

var spaceConnections = make(map[*websocket.Conn]SpaceConnection)
var boardConnections = make(map[*websocket.Conn]BoardConnection)

func sendSpaceMessage(message interface{}, spaceID uint) {
	for conn, connData := range spaceConnections {
		if connData.SpaceID == spaceID {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("error: %v", err)
				conn.Close()
				delete(spaceConnections, conn)
			}
		}
	}
}

func sendBoardMessage(message interface{}, spaceID, boardID uint) {
	for conn, connData := range boardConnections {
		if connData.SpaceID == spaceID && connData.BoardID == boardID {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("error: %v", err)
				conn.Close()
				delete(boardConnections, conn)
			}
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Получаем ID пространства и доски из заголовков или параметров запроса
	spaceID := 0 // Получаем ID пространства
	boardID := 0 // Получаем ID доски (может быть 0, если соединение для пространства)

	if boardID == 0 {
		// Соединение для пространства
		spaceConnections[conn] = SpaceConnection{SpaceID: spaceID, Conn: conn}
	} else {
		// Соединение для доски
		boardConnections[conn] = BoardConnection{SpaceID: spaceID, BoardID: boardID, Conn: conn}
	}
}
*/
