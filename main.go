package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type Shape struct {
	Triangle string `json:"triangle"`
	Square   string `json:"square"`
	Circle   string `json:"circle"`
}

var shapes = Shape{}

func InsertData(db *sql.DB,data string){
  dataFormat := fmt.Sprintf("INSERT INTO shape (%s) VALUES (?)",data)
  stmt, err := db.Prepare(dataFormat)
  if err != nil {
    panic(err)
  }
  defer stmt.Close()

  _, err = stmt.Exec(1)
  if err != nil{
    panic(err)
  }
}
func main() {

	//connect database
	db, err := sql.Open("sqlite3", "mydb.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// Upgrader is used to upgrade http requests to websocket connections
	upgrader := websocket.Upgrader{
		// Read and write buffer size
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// allow origin
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// HandleFunc is used to define a websocket endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the http request to a websocket connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Listen for messages from the client
		for {
      msgType, data, err := ws.ReadMessage()
      if err != nil {
        fmt.Println(err)
        break
      }
      InsertData(db,string(data))
			rows, err := db.Query("SELECT * FROM shape")
			if err != nil {
				panic(err)
			}
			// Iterate over the rows and print the results.
			for rows.Next() {
				err := rows.Scan(&shapes.Triangle, &shapes.Square, &shapes.Circle)
				if err != nil {
					panic(err)
				}

			}

			jsonBytes, err := json.Marshal(shapes)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Send the message back to the client
			ws.WriteMessage(msgType, jsonBytes)
		}
	})

	// Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
