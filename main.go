package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const pages = 2

var (
	config     *oauth1.Config
	token      *oauth1.Token
	httpClient *http.Client
	client     *twitter.Client
)

func init() {
	config = oauth1.NewConfig(os.Getenv("APIKEY"), os.Getenv("APISECRET"))
	token = oauth1.NewToken(os.Getenv("TOKEN"), os.Getenv("TOKENSECRET"))
	httpClient = config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:           "127.0.0.1:" + os.Getenv("PORT"),
		Handler:        makeMuxRouter(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Kudurru listening on port:", os.Getenv("PORT"))
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleRoot).Methods("GET")
	muxRouter.HandleFunc("/home", handleHomeTimeline).Methods("GET")
	muxRouter.HandleFunc("/user/{id}", handleUserTimeline).Methods("GET")
	return muxRouter
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Kudurru, written in Stone 🗿")
}

func handleHomeTimeline(w http.ResponseWriter, r *http.Request) {
	tweets, resp, err := client.Timelines.HomeTimeline(
		&twitter.HomeTimelineParams{Count: 10},
	)
	if err != nil {
		respondWithError(err, w)
	}
	defer resp.Body.Close()

	for _, tweet := range tweets {
		w.Write([]byte(tweet.Text))
	}
}

func handleUserTimeline(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	userHandle := muxVars["id"]

	tweets, resp, err := client.Timelines.UserTimeline(
		&twitter.UserTimelineParams{
			ScreenName:     userHandle,
			Count:          10,
			TrimUser:       twitter.Bool(false),
			ExcludeReplies: twitter.Bool(true),
			TweetMode:      "extended"})
	if err != nil {
		respondWithError(err, w)
	}
	defer resp.Body.Close()

	for _, tweet := range tweets {
		w.Write([]byte(tweet.FullText))
		fmt.Fprintf(w, "\n>\n")
	}
}

func respondWithError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
