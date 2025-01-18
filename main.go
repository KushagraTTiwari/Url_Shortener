package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalUrl  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

/*
   ShortUrl --> {
                   ID
                   OriginalUrl
                   ShortUrl
                   CreationData
               }
*/

var urlDB = make(map[string]URL)

func generateShortURL(OriginalUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalUrl)) //It converts the originalURL string to a byte slice
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func createURL(OriginalUrl string) string {
	shortUrl := generateShortURL(OriginalUrl)
	id := shortUrl // Use the short URL as the ID for simplicity
	urlDB[id] = URL{
		ID:           id,
		OriginalUrl:  OriginalUrl,
		ShortUrl:     shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invailid request body", http.StatusBadRequest)
		return
	}

	shortURL := createURL(data.URL)
	// fmt.Fprintf(w, shortURL)
	response := struct {
		ShortedURL string `json:"short_url"`
	}{ShortedURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):] // it will take the whole string after the "redirect"
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invailid Request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}

func main() {
	fmt.Println("Url-Shortner")

	//Register the handler function to handle all requests to the root URL ("/")
	http.HandleFunc("/", handler)
	http.HandleFunc("/short", ShortURLHandler)
	http.HandleFunc("/redirect/", RedirectURLHandler)

	// start the http server on port 3000
	fmt.Println("Starting server on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error on starting server:", err)
	}
}
