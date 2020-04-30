package kudurru

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// StartupMessage when the server fires
const StartupMessage = "Kudurru, written in Stone 🗿"

var (
	config     *clientcredentials.Config
	httpClient *http.Client
	twClient   *twitter.Client
)

// New creates the server to connect with Twitter
func New() *http.Server {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ipfsGWY = shell.NewShell("https://ipfs.infura.io:5001")

	config = &clientcredentials.Config{
		ClientID:     os.Getenv("APIKEY"),
		ClientSecret: os.Getenv("APISECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient = config.Client(oauth2.NoContext)
	twClient = twitter.NewClient(httpClient)

	return &http.Server{
		Addr:           "127.0.0.1:" + os.Getenv("PORT"),
		Handler:        makeMuxRouter(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", HandleRoot).Methods("GET")
	muxRouter.HandleFunc("/home", HandleHomeTimeline).Methods("GET")
	muxRouter.HandleFunc("/user/{id}", HandleUserTimeline).Methods("GET")
	return muxRouter
}

// HandleRoot serves the / path
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, StartupMessage)
}

// HandleHomeTimeline fetches the authenticated home time line
func HandleHomeTimeline(w http.ResponseWriter, r *http.Request) {
	tweets, resp, err := twClient.Timelines.HomeTimeline(
		&twitter.HomeTimelineParams{Count: 11},
	)
	if err != nil {
		respondWithError(err, w)
	}
	defer resp.Body.Close()

	for _, tweet := range tweets {
		w.Write([]byte(tweet.Text))
	}
}

// HandleUserTimeline fetches the time line of the provided user
func HandleUserTimeline(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	userHandle := muxVars["id"]

	tweets, resp, err := twClient.Timelines.UserTimeline(
		&twitter.UserTimelineParams{
			ScreenName:     userHandle,
			Count:          11,
			TrimUser:       twitter.Bool(false),
			ExcludeReplies: twitter.Bool(true),
			TweetMode:      "extended"})
	if err != nil {
		respondWithError(err, w)
	}
	defer resp.Body.Close()

	for _, tweet := range tweets {
		var buffer strings.Builder
		buffer.WriteString(tweet.FullText)
		buffer.WriteString("\n\n")
		buffer.WriteString(strings.Repeat("+", 20))
		buffer.WriteString("\n\n")
		w.Write([]byte(buffer.String()))
	}
}

func respondWithError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
