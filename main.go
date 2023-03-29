package main

import (
	"fmt"

	"net/http"
	"os"

	"example.com/games/controllers"
	"example.com/games/database"
	"example.com/games/environment"
	_ "github.com/go-sql-driver/mysql"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	//"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Val  string `json:"val"`
	Room string `json:"room"`
}

// var clients = make(map[*websocket.Conn]bool)
// var broadcaster = make(chan ChatMessage)

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// ensure connection close when function returns
// 	defer ws.Close()
// 	clients[ws] = true

// 	for {
// 		var msg ChatMessage
// 		// Read in a new message as JSON and map it to a Message object
// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			delete(clients, ws)
// 			break
// 		}
// 		// send new message to the channel
// 		broadcaster <- msg
// 	}
// }

// // If a message is sent while a client is closing, ignore the error
// func unsafeError(err error) bool {
// 	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
// }

// func handleMessages() {
// 	for {
// 		// grab any next message from channel
// 		msg := <-broadcaster

// 		messageClients(msg)
// 	}
// }

// func messageClients(msg ChatMessage) {
// 	// send to every client currently connected
// 	for client := range clients {
// 		messageClient(client, msg)
// 	}
// }

//	func messageClient(client *websocket.Conn, msg ChatMessage) {
//		err := client.WriteJSON(msg)
//		if err != nil && unsafeError(err) {
//			log.Printf("error: %v", err)
//			client.Close()
//			delete(clients, client)
//		}
//	}
func main() {
	db := os.Getenv("DATABASE")
	if db == "" {
		db = environment.ViperEnvVariable("DATABASE")
	}
	client, ctx, _, err, _ := database.Connect(db)
	if err != nil {
		panic(err)
	}

	database.Ping(client, ctx)
	r := mux.NewRouter()
	r.HandleFunc("/api/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("/api/login", controllers.Login).Methods("POST")
	r.HandleFunc("/{room}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})
	server := socketio.NewServer(nil)

	server.OnEvent("/", "join-room", func(s socketio.Conn, msg string) string {
		fmt.Println("->", msg)
		if msg == "game-started" {
			return "game-started"
		}
		s.Join(msg)
		return "joined"
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, msg ChatMessage) string {
		fmt.Println(msg)
		server.BroadcastToRoom("/", msg.Room, "room-chat", msg.Val)

		return "you: " + msg.Val
	})

	server.OnEvent("/", "remove-socket", func(s socketio.Conn, room string) {
		s.Leave(room)
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	r.Handle("/socket.io/", server)

	port := os.Getenv("PORT")
	if port == "" {
		port = environment.ViperEnvVariable("PORT")
	}
	fmt.Println("running on port: ", port)

	http.ListenAndServe(":"+port, r)
}

// package main

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// func main() {
// 	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
// 		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

// 		for {
// 			// Read message from browser
// 			msgType, msg, err := conn.ReadMessage()
// 			if err != nil {
// 				return
// 			}

// 			// Print the message to the console
// 			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

// 			// Write message back to browser
// 			if err = conn.WriteMessage(msgType, msg); err != nil {
// 				return
// 			}
// 		}
// 	})

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("ent")
// 		http.ServeFile(w, r, "./public/index.html")
// 	})

// 	http.ListenAndServe(":8080", nil)
// }
