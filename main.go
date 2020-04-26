package main

import "github.com/dghubble/oauth1"

// Tweet that represents tweets from twitter
type Tweet struct {
	Date string `json:"created_at"`
	Text string `json:"text"`
	ID   string `json:"id_str"`
}

var config *oauth1.Config
var token *oauth1.Token

const pages = 18
