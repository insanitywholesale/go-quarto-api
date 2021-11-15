package main

import (
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
	"log"
	"net/http"
	"os"
)

var (
//commit date + hash vars
)

//Variable of all Quarto pieces
var AllQuartoPieces = [16]*QuartoPiece{
	//All false
	&QuartoPiece{
		Dark:   false,
		Short:  false,
		Hollow: false,
		Round:  false,
	},
	//One true
	&QuartoPiece{
		Dark:   true,
		Short:  false,
		Hollow: false,
		Round:  false,
	},
	&QuartoPiece{
		Dark:   false,
		Short:  true,
		Hollow: false,
		Round:  false,
	},
	&QuartoPiece{
		Dark:   false,
		Short:  false,
		Hollow: true,
		Round:  false,
	},
	&QuartoPiece{
		Dark:   false,
		Short:  false,
		Hollow: false,
		Round:  true,
	},
	//Two true
	&QuartoPiece{
		Dark:   true,
		Short:  true,
		Hollow: false,
		Round:  false,
	},
	&QuartoPiece{
		Dark:   false,
		Short:  true,
		Hollow: true,
		Round:  false,
	},
	&QuartoPiece{
		Dark:   false,
		Short:  true,
		Hollow: false,
		Round:  true,
	},
	&QuartoPiece{
		Dark:   true,
		Short:  false,
		Hollow: true,
		Round:  false,
	},
	&QuartoPiece{
		Dark:   true,
		Short:  false,
		Hollow: false,
		Round:  true,
	},
	&QuartoPiece{
		Dark:   false,
		Short:  false,
		Hollow: true,
		Round:  true,
	},
	//Three true
	&QuartoPiece{
		Dark:   false,
		Short:  true,
		Hollow: true,
		Round:  true,
	},
	&QuartoPiece{
		Dark:   true,
		Short:  false,
		Hollow: true,
		Round:  true,
	},
	&QuartoPiece{
		Dark:   true,
		Short:  true,
		Hollow: false,
		Round:  true,
	},
	&QuartoPiece{
		Dark:   true,
		Short:  true,
		Hollow: true,
		Round:  false,
	},
	//All true
	&QuartoPiece{
		Dark:   true,
		Short:  true,
		Hollow: true,
		Round:  true,
	},
}

// Variable of empty game board
var EmptyBoard = [4][4]*QuartoPiece{
	{&QuartoPiece{},&QuartoPiece{},&QuartoPiece{},&QuartoPiece{}},
	{&QuartoPiece{},&QuartoPiece{},&QuartoPiece{},&QuartoPiece{}},
	{&QuartoPiece{},&QuartoPiece{},&QuartoPiece{},&QuartoPiece{}},
	{&QuartoPiece{},&QuartoPiece{},&QuartoPiece{},&QuartoPiece{}},
}

// Constant for Bad Request
const BadReq string = `{"error": "bad request"}`

// Constant for Not Found
const NotFound string = `{"error": "not found"}`

// Constant for User Not Found
const UserNotFound string = `{"error": "user not found"}`

// Constant for Game Not Found
const GameNotFound string = `{"error": "game not found"}`

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
	GameId         string     `json:"game_id"`
	ActivePlayers  []*UserId  `json:"players"`
	InvitedPlayers []*UserId  `json:"invited_players"`
	PlayerTurn     *UserId    `json:"next_player"`
	ActivityStatus bool       `json:"activity_status"`
	State          *GameState `json:"game_state"`
}

//TODO: fill in with fields
type GameState struct {
	Board        [4][4]*QuartoPiece
	UnusedPieces [16]*QuartoPiece
}

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
	//TODO: replace with db call
	testUsers = append(testUsers, u)
	uid := &UserId{
		UserName: u.UserName,
		UserId:   shortid.MustGenerate(),
	}
	//TODO: replace with db call
	testUserIds = append(testUserIds, uid)
	json.NewEncoder(w).Encode(uid)
}

func getGameState(w http.ResponseWriter, r *http.Request) {
	log.Println("getGameState called")
	w.Header().Set("Content-Type", "application/json")
	//get the path parameters
	params := mux.Vars(r)
	//get game_id from path param
	gameId, _ := params["game_id"]
	for _, g := range testGames {
		if g.GameId == gameId {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(g.State)
			return
		}
	}
	//TODO: errors if game not found, user not authorized
}

func createGame(w http.ResponseWriter, r *http.Request) {
	log.Println("createGame called")
	w.Header().Set("Content-Type", "application/json")
	//user that creates the game
	uid := &UserId{}
	err := json.NewDecoder(r.Body).Decode(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(BadReq))
		return
	}
	//create a new game instance
	g := &Game{
		GameId:         shortid.MustGenerate(),
		ActivityStatus: true,
		State: &GameState{
			Board: EmptyBoard,
			UnusedPieces: AllQuartoPieces,
		},
	}
	//automatically invite the game creator to the game
	g.InvitedPlayers = append(g.InvitedPlayers, uid)
	//TODO: replace with db call
	testGames = append(testGames, g)
	//TODO: generate player-specific game code to return
	json.NewEncoder(w).Encode(g)
}

func inviteToGame(w http.ResponseWriter, r *http.Request) {
	log.Println("inviteToGame called")
	w.Header().Set("Content-Type", "application/json")
	//get the path parameters
	params := mux.Vars(r)
	//get game_id from path param
	gameId, _ := params["game_id"]
	log.Println("gameId", gameId)

	//user to be invited
	var uid *UserId = nil
	//get the name of the user to be invited from path param
	inviteeName, _ := params["username"]
	log.Println("inviteeName", inviteeName)
	//see if user exists in the user database
	for idx, u := range testUserIds {
		if u.UserName == inviteeName {
			uid = testUserIds[idx]
			break
		}
	}
	//return error if user with username can't be found
	if uid == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(GameNotFound))
		return
	}
	//append player to game if game exists
	for _, g := range testGames {
		if g.GameId == gameId {
			g.InvitedPlayers = append(g.InvitedPlayers, uid)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
			return
		}
	}
	//return error if game doesn't exist
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(GameNotFound))
	return
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	log.Println("joinGame called")
	w.Header().Set("Content-Type", "application/json")
	//get the path parameters
	params := mux.Vars(r)
	//get game_id from path param
	gameId, _ := params["game_id"]
	log.Println("gameId", gameId)

	//user trying to join
	uid := &UserId{}
	err := json.NewDecoder(r.Body).Decode(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(BadReq))
		return
	}

	//TODO: error messages in following loop
	//TODO: returns in following loop
	for _, g := range testGames {
		if g.GameId == gameId {
			for _, u := range g.InvitedPlayers {
				if cap(g.ActivePlayers) <= len(g.ActivePlayers) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error": "couldn't join because lobby is full"}`))
					return
				}
				if uid.UserId == u.UserId {
					g.ActivePlayers = append(g.ActivePlayers, uid)
					g.InvitedPlayers = g.InvitedPlayers[:len(g.InvitedPlayers)-1]
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"message": "success"}`))
					return
				}
			}
		}
	}
}

func playInGame(w http.ResponseWriter, r *http.Request) {
	log.Println("playInGame called")
	w.Header().Set("Content-Type", "application/json")
	//TODO: depends on game, will probably be quarto
}

func ifQuarto(qp [4]*QuartoPiece) bool {
	//TODO: replace with x, y, z, w vars for readability
	if qp[0].Dark == qp[1].Dark == qp[2].Dark == qp[3].Dark {
		return true
	} else if qp[0].Short == qp[1].Short == qp[2].Short == qp[3].Short {
		return true
	} else if qp[0].Hollow == qp[1].Hollow == qp[2].Hollow == qp[3].Hollow {
		return true
	} else if qp[0].Round == qp[1].Round == qp[2].Round == qp[3].Round {
		return true
	} else {
		return false
	}
}

func checkGameState(gameId string) bool {
	var gameState *GameState
	for _, g := range testGames {
		if g.GameId == gameId {
			gameState = g.State
		}
	}
	board := gameState.Board
	unusedPieces := gameState.UnusedPieces
	log.Println(unusedPieces)
	diag1 := [4]*QuartoPiece{board[0][0], board[1][1], board[2][2], board[3][3]}
	diag2 := [4]*QuartoPiece{board[0][3], board[1][2], board[2][1], board[3][0]}
	for i, row := range board {
		log.Println(i, row)
		// Don't bother if 4 pieces haven't been on the board
		if cap(unusedPieces) > 12 {
			break
		}
		// Don't bother if row isn't full
		if cap(row) < 4 {
			break
		}
		if ifQuarto(row) {
			return true
		}
		var col [4]*QuartoPiece
		for j, colItem := range row {
			log.Println(j, col)
			log.Println(j, colItem)
			col[j] = colItem
		}
		if cap(col) == 4 && ifQuarto(col) {
			return true
		}
		if ifQuarto(diag1) {
			return true
		}
		if ifQuarto(diag2) {
			return true
		}
	}
	return false
}

// Function to set server HTTP port
func setupHTTPPort() string {
	httpPort := os.Getenv("QUARTO_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8000"
	}
	return httpPort
}

// Only for testing
func genRandomPiece() *QuartoPiece {
	qp := &QuartoPiece{
		Dark:   randomdata.Boolean(),
		Short:  randomdata.Boolean(),
		Hollow: randomdata.Boolean(),
		Round:  randomdata.Boolean(),
	}
	return qp
}

func setupRouter() http.Handler {
	// Set up router
	router := mux.NewRouter()
	// Set up subrouter for user functions
	userRouter := router.PathPrefix("/user").Subrouter()
	// Set up subrouter for game functions
	gameRouter := router.PathPrefix("/game").Subrouter()
	// Set up routes for user API
	userRouter.HandleFunc("/register", createUser)
	// Set up routes for game API
	//gameRouter.HandleFunc("/new", createGame)
	gameRouter.HandleFunc("", createGame)
	gameRouter.HandleFunc("/{game_id}", getGameState)
	gameRouter.HandleFunc("/{game_id}/join", joinGame)
	gameRouter.HandleFunc("/{game_id}/play", playInGame)
	gameRouter.HandleFunc("/{game_id}/invite/{username}", inviteToGame)
	return router
}

func main() {
	// Determine port to run at
	httpPort := setupHTTPPort()
	// Set up the router for the API
	router := setupRouter()
	// Print a message so there is feedback to the app admin
	log.Println("starting server at port", httpPort)
	// One-liner to start the server or print error
	log.Fatal(http.ListenAndServe(":"+httpPort, router))
}
