package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
	"log"
	"net/http"
	"os"
)

var (
	httpPort string
)

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

//TODO: rethink active/inactive players thing
type Game struct {
	GameId         string    `json:"game_id"`
	ActivePlayers  []*UserId `json:"players"`
	InvitedPlayers []*UserId `json:"invited_players"`
	ActivityStatus bool      `json:"activity_status"`
	State          GameState `json:"game_state"`
}

//TODO: fill in with fields
type GameState struct{}

type QuartoPiece struct {
	Dark   bool
	Short  bool
	Hollow bool
	Round  bool
}

//TODO: replace with database(s)
var testUsers []*User
var testUserIds []*UserId
var testPlayers []*UserId
var testGames []*Game

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
	//TODO: save uid to storage
	json.NewEncoder(w).Encode(uid)
}

func getGameState(w http.ResponseWriter, r *http.Request) {
	log.Println("createUser called")
	w.Header().Set("Content-Type", "application/json")
	gameState := &GameState{}
	//TODO: for gamelist if game's game_id equals provided game_id
	//assign game.State to gameState
	json.NewEncoder(w).Encode(gameState)
}

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
	g.InvitedPlayers = append(g.InvitedPlayers, uid)
	//TODO: generate player-specific game code to return
	json.NewEncoder(w).Encode(g)
}

func inviteToGame(w http.ResponseWriter, r *http.Request) {
	//TODO: for each game's game_id if it equals provided game_id
	//game.InvitedPlayers append requesting player so they can join
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	//TODO: for gamelist if game's game_id equals provided game_id
	//game.ActivePlayers append requesting player so they can join
}

func playInGame(w http.ResponseWriter, r *http.Request) {
	//TODO: depends on game, will probably be quarto
}

// Function to set server HTTP port
func setupHTTPPort() {
	httpPort = os.Getenv("QUARTO_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
}

func main() {
	// Set up router
	router := mux.NewRouter()
	// Set up subrouter for user functions
	userRouter := router.PathPrefix("/user").Subrouter()
	// Set up subrouter for game functions
	gameRouter := router.PathPrefix("/game").Subrouter()
	// Set up routes for user API
	userRouter.HandleFunc("/register", createUser)
	// Set up routes for game API
	gameRouter.HandleFunc("/new", createGame)
	gameRouter.HandleFunc("/{game_id}", getGameState)
	gameRouter.HandleFunc("/{game_id}/join", joinGame)
	gameRouter.HandleFunc("/{game_id}/play", playInGame)
	gameRouter.HandleFunc("/{game_id}/invite/{username}", inviteToGame)
	// Determine port to run at
	httpPort()
	// Print a message so there is feedback to the app admin
	log.Println("starting server at port" + httpPort)
	// One-liner to start the server or print error
	log.Fatal(http.ListenAndServe(":"+httpPort, router))
}
