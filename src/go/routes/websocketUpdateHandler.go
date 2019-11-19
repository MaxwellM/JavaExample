package routes

import (
	"Orbit-Messenger/src/go/db"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func (rc *RouteController) UpdateHandler(wsConn *websocket.Conn, state *State) {
	serverActionLen := rc.serverActions.ActionCount
	for !state.LoggedIn {
		// waits for the user to login
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("user logged in")

	for {
		// checks if the user has logged out
		if state.LoggedOut {
			log.Println("user logged out")
			return
		}
		if serverActionLen != rc.serverActions.ActionCount {
			clientAction, err := rc.serverActions.GetNewestAction()
			if err != nil {
				wsConn.WriteJSON(err.Error())
			} else {
				data := map[string]int64{
					clientAction.Action: clientAction.messageId,
				}
				wsConn.WriteJSON(data)
				serverActionLen = rc.serverActions.ActionCount
			}
		}

		var serverResponse ServerResponse
		if state.LastMessageId == 0 {
			serverResponse = rc.getAllMessagesForClient(&state.LastMessageId)
		} else {
			serverResponse = rc.getNewMessages(&state.LastMessageId)
		}

		if serverResponse.Messages != nil {
			fmt.Println("sending message")
			wsConn.WriteJSON(serverResponse)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (rc RouteController) getNewMessages(lastMessageId *int64) ServerResponse {
	var serverResponse ServerResponse
	messages, err := rc.dbConn.GetNewestMessagesFrom(*lastMessageId)

	if err != nil {
		serverResponse.Errors = err.Error()
		return serverResponse
	}
	if len(messages) == 0 {
		return serverResponse
	}

	serverResponse.Messages = messages
	updateLastMessageId(messages, lastMessageId)
	return serverResponse
}

// Gets all the messages for the client
func (rc RouteController) getAllMessagesForClient(lastMessageId *int64) ServerResponse {
	log.Println("getting All Messages")
	var serverResponse ServerResponse
	activeUsers, err := rc.dbConn.GetUsersByStatus(true)

	// send error to client
	if err != nil {
		serverResponse.Errors = err.Error()
		return serverResponse
	}
	messages, err := rc.dbConn.GetAllMessages()
	if err != nil {
		serverResponse.Errors = err.Error()
		return serverResponse
	}

	rooms, err := rc.dbConn.GetAllRooms()


	serverResponse.ActiveUsers = activeUsers
	serverResponse.Messages = messages
	serverResponse.Rooms = rooms
	updateLastMessageId(messages, lastMessageId)
	return serverResponse
}

func updateLastMessageId(messages []db.Message, lastMessageId *int64) {
	if len(messages) < 1 {
		return
	}
	*lastMessageId = messages[len(messages)-1].MessageId
}
