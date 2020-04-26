package main

import (
	"log"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
)

// Tweet that represents tweets from twitter
type Tweet struct {
	Date string `json:"created_at"`
	Text string `json:"text"`
	ID   string `json:"id_str"`
}

var config *oauth1.Config
var token *oauth1.Token

const pages = 1

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/{id}", handleGetTweets).Methods("GET")
	return muxRouter
}

func handleGetTweets(w http.ResponseWriter, r *http.Request) {

}

func respondWithError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
