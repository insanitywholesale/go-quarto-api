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
	// variables for commit hash and commit date
	commitHash string
	commitDate string
)

// Variable of all Quarto pieces
var AllQuartoPieces = [16]*QuartoPiece{
	// All false
	&QuartoPiece{
		Dark:   false,
		Short:  false,
		Hollow: false,
		Round:  false,
	},
	// One true
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
	// Two true
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
	// Three true
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
	// All true
	&QuartoPiece{
		Dark:   true,
		Short:  true,
		Hollow: true,
		Round:  true,
	},
}

// Variable of empty game board
var EmptyBoard = [4][4]*QuartoPiece{
	{&QuartoPiece{}, &QuartoPiece{}, &QuartoPiece{}, &QuartoPiece{}},
	{&QuartoPiece{}, &QuartoPiece{}, &QuartoPiece{}, &QuartoPiece{}},
	{&QuartoPiece{}, &QuartoPiece{}, &QuartoPiece{}, &QuartoPiece{}},
	{&QuartoPiece{}, &QuartoPiece{}, &QuartoPiece{}, &QuartoPiece{}},
}

// Constant for maximum amount of players per game
const MaxPlayers int = 2

// Constant for Bad Request
const BadReq string = `{"error": "bad request"}`

// Constant for Not Found
const NotFound string = `{"error": "not found"}`

// Constant for Unauthorized
const Unauth string = `{"error": "unauthorized"}`

// Constant for Unauthorized
const ServerError string = `{"error": "internal server error"}`

// Constant for success message
const MsgSuccess string = `{"message": "success"}`

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
	GameId         string       `json:"game_id"`
	ActivePlayers  []*UserId    `json:"active_players"`
	InvitedPlayers []*UserId    `json:"invited_players"`
	NextPlayer     *UserId      `json:"next_player"` //TODO: move to GameState
	NextPiece      *QuartoPiece `json:"next_piece"`  //TODO: move to GameState
	ActivityStatus bool         `json:"activity_status"`
	State          *GameState   `json:"game_state"`
	Winner         *UserId      `json:"winner"`
}

//TODO: fill in with fields
type GameState struct {
	Board        [4][4]*QuartoPiece `json:"board"`
	UnusedPieces [16]*QuartoPiece   `json:"unused_pieces"`
}

type GameMove struct {
	PositionX int32        `json:"position_x"`
	PositionY int32        `json:"position_y"`
	NextPiece *QuartoPiece `json:"next_piece"`
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
			Board:        EmptyBoard,
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
			w.Write([]byte(MsgSuccess))
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
	//TODO: unindent following loop
	for _, g := range testGames {
		if g.GameId == gameId {
			for _, u := range g.InvitedPlayers {
				if cap(g.ActivePlayers) == MaxPlayers {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error": "couldn't join because game is full"}`))
					return
				} else if cap(g.ActivePlayers) > MaxPlayers {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error": "I honestly don't know how this happened"}`))
					return
				}
				if uid.UserId == u.UserId {
					g.ActivePlayers = append(g.ActivePlayers, uid)
					g.InvitedPlayers = g.InvitedPlayers[:len(g.InvitedPlayers)-1]
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(MsgSuccess))
					return
				}
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(GameNotFound))
	return
}

func playInGame(w http.ResponseWriter, r *http.Request) {
	log.Println("playInGame called")
	w.Header().Set("Content-Type", "application/json")
	for _, g := range testGames {
		if g.GameId == gameId {
			if len(g.ActivePlayers) != 2 {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "need exactly 2 players to play this game"}`))
				return
			}
			// requesting player
			u := &UserId{}
			err := json.NewDecoder(r.Body).Decode(uid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(BadReq))
				return
			}
			// player playing next
			player := g.NextPlayer
			// if requesting player s not player playing next, error out
			if player.UserId != u.UserId {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(Unauth))
				return
			}
			// piece to be placed
			var piece *QuartoPiece
			// make sure the game's next piece has been set
			if g.NextPiece == nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(ServerError))
				return
			} else {
				piece = g.NextPiece
			}
			// game move
			gameMove := &GameMove{}
			err := json.NewDecoder(r.Body).Decode(gameMove)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(BadReq))
				return
			}
			g.Board[gameMove.PositionX][gameMove.PositionY] = g.NextPiece
			// make sure the game move's next piece has been set
			if gameMove.NextPiece == nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(ServerError))
				return
			} else {
				g.NextPiece = gameMove.NextPiece
			}
			//TODO: check if quatro and such
			//TODO: deal with ActivityStatus
			done := checkGameState(g.GameId)
			if done {
				g.ActivityStatus = false
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "` + u.UserId + ` is the winner!"}`))
				return
			}
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(GameNotFound))
	return
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
	log.Println("unusedPieces", unusedPieces)
	// Statically define diagonal and reverse diagonal
	diag1 := [4]*QuartoPiece{board[0][0], board[1][1], board[2][2], board[3][3]}
	diag2 := [4]*QuartoPiece{board[0][3], board[1][2], board[2][1], board[3][0]}
	// Go through the board and check if anything qualifies as quarto
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
		// Check if current row has quarto
		if ifQuarto(row) {
			return true
		}
		// Collect items from column
		var col [4]*QuartoPiece
		for j, colItem := range row {
			log.Println(j, col)
			log.Println(j, colItem)
			col[j] = colItem
		}
		// Check if there are 4 pieces in the column
		if cap(col) == 4 && ifQuarto(col) {
			return true
		}
		// Check if there are 4 pieces in the diagonal
		if cap(diag1) == 4 && ifQuarto(diag1) {
			return true
		}
		// Check if there are 4 pieces in the reverse diagonal
		if cap(diag2) == 4 && ifQuarto(diag2) {
			return true
		}
	}
	// Return false if none of the above succeeded
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

// Generate a random piece
func genRandomPiece() *QuartoPiece {
	return &QuartoPiece{
		Dark:   randomdata.Boolean(),
		Short:  randomdata.Boolean(),
		Hollow: randomdata.Boolean(),
		Round:  randomdata.Boolean(),
	}
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
