package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
	"log"
	"net/http"
)

// lifted from go tutorial
/*
func main() {
	// Set up router
	router := mux.NewRouter()
	// Set up subrouter for api version 1
	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	// Set up routes
	apiV1.HandleFunc("/deliveries", GetAllDeliveries).Methods(http.MethodGet)
	// Start http server
	log.Fatal(http.ListenAndServe(":8000", router))
}
*/
// Constant for Bad Request
const BadReq string = `{"error": "bad request"}`

// Constant for Not Found
const NotFound string = `{"error": "not found"}`

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserId struct {
	UserName string `json:"username"`
	UserId   string `json:"user_id"`
}

type Game struct {
	GameId         string    `json:"game_id"`
	Players        []*UserId `json:"players"`
	ActivityStatus bool      `json:"activity_status"`
	State          GameState `json:"game_state"`
}

type GameState struct{}

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("createUser called")
	w.Header().Set("Content-Type", "application/json")
	u := &User{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(BadReq))
		return
	}
	uid := &UserId{
		UserName: u.UserName,
		UserId:   shortid.MustGenerate(),
	}
	//save uid to storage
	json.NewEncoder(w).Encode(uid)
}

func getGameState(w http.ResponseWriter, r *http.Request) {}

func createGame(w http.ResponseWriter, r *http.Request) {
	log.Println("createGame called")
	w.Header().Set("Content-Type", "application/json")
	uid := &UserId{}
	err := json.NewDecoder(r.Body).Decode(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(BadReq))
		return
	}
	g := &Game{
		GameId:         shortid.MustGenerate(),
		ActivityStatus: true,
	}
	g.Players = append(g.Players, uid)
	//generate player-specific game code to return
	json.NewEncoder(w).Encode(g)
}

func inviteToGame(w http.ResponseWriter, r *http.Request) {}
func joinGame(w http.ResponseWriter, r *http.Request)     {}
func playInGame(w http.ResponseWriter, r *http.Request)   {}

func main() {
	// Set up router
	router := mux.NewRouter()
	router.HandleFunc("/user/register", createUser)
	router.HandleFunc("/game/new", createGame)
	router.HandleFunc("/game/{game_id}", getGameState)
	router.HandleFunc("/game/invite/{username}", inviteToGame)
	router.HandleFunc("/game/{game_id}/join", joinGame)
	router.HandleFunc("/game/{game_id}/play", playInGame)
	log.Println("starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
