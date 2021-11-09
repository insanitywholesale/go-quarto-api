package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

func createUser(w http.ResponseWriter, r *http.Request) {}

func getGameState(w http.ResponseWriter, r *http.Request) {}
func inviteToGame(w http.ResponseWriter, r *http.Request) {}
func joinGame(w http.ResponseWriter, r *http.Request) {}
func playInGame(w http.ResponseWriter, r *http.Request) {}

func main() {
	// Set up router
	router := mux.NewRouter()
	router.HandleFunc("/user/register", createUser)
	router.HandleFunc("/game/{game_id}", getGameState)
	router.HandleFunc("/game/new", createGame)
	router.HandleFunc("/game/invite/{username}", inviteToGame)
	router.HandleFunc("/game/{game_id}/join", joinGame)
	router.HandleFunc("/game/{game_id}/play", playInGame)
}
