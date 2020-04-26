package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Tweet that represents tweets from twitter
type Tweet struct {
	Date string `json:"created_at"`
	Text string `json:"text"`
	ID   string `json:"id_str"`
}

var config *oauth1.Config
var token *oauth1.Token

const pages = 2

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleHome).Methods("GET")
	muxRouter.HandleFunc("/{id}", handleGetTweets).Methods("GET")
	return muxRouter
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Kudurru, written in Stone 🗿")
}

func handleGetTweets(w http.ResponseWriter, r *http.Request) {
	var maxIDQuery string
	// var tweets []Tweet
	muxVars := mux.Vars(r)
	userHandle := muxVars["id"]

	httpClient := config.Client(oauth1.NoContext, token)

	for i := 0; i < pages; i++ {
		path := fmt.Sprintf("https://api.twitter.com/1.1/statues/user_timeline.json?screen_name=%v&include_rts=false&count=10%v", userHandle, maxIDQuery)

		if strings.Contains(path, "favicon.ico") {
			break
		}

		resp, err := httpClient.Get(path)
		if err != nil {
			respondWithError(err, w)
			break
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			respondWithError(err, w)
			break
		}

		var gotTweets []Tweet
		err = json.Unmarshal(body, &gotTweets)
		if err != nil {
			respondWithError(err, w)
			break
		}

		fmt.Println(gotTweets)
	}
}

func respondWithError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	config = oauth1.NewConfig(os.Getenv("APIKEY"), os.Getenv("APISECRET"))
	token = oauth1.NewToken(os.Getenv("TOKEN"), os.Getenv("TOKENSECRET"))

	s := &http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        makeMuxRouter(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
