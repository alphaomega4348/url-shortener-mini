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
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// Storing URL in RunTime Memory can be scaled to a database
var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	//md5 hash of original url
	hasher := md5.New()
	//write original url to hasher
	hasher.Write([]byte(OriginalURL))
	//get the hash as a byte array
	data := hasher.Sum(nil)
	//convert byte array to hex string
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func generateAndStoreURLtoDB(OriginalURL string) string {
	shortURL := generateShortURL(OriginalURL)
	//Storing in Map
	urlDB[shortURL] = URL{
		ID:           shortURL,
		OriginalURL:  OriginalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURLfromDB(id string) (URL, error) {
	url, isValid := urlDB[id]
	if isValid {
		return url, nil
	}
	return url, errors.New("URL not found")
}

func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	var receivedData struct {
		URL string `json:"url"`
	}
	//get the url from the request and store it in a variable
	err := json.NewDecoder(r.Body).Decode(&receivedData)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	//Create Short URL
	shortURL := generateAndStoreURLtoDB(receivedData.URL)
	//fmt.Fprintf(w, shortURL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

}

func redirectToOriginalURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect"):]
	url, err := getURLfromDB(id)
	if err != nil {
		http.Error(w, "URL Not Found", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)

}

func HomePageURL(w http.ResponseWriter, r *http.Request) {
	//Fprintf helps to write to the response writer
	fmt.Fprintf(w, "Get URL")
}

func main() {
	fmt.Println("Starting")
	shortened := generateAndStoreURLtoDB("https://www.google.com")
	fmt.Println(shortened)

	//Handler function to handle all requests to home route "/"
	http.HandleFunc("/", HomePageURL)
	http.HandleFunc("/shorten", shortURLHandler)
	http.HandleFunc("/redirect", redirectToOriginalURLHandler)

	//Creating a new HTTP server on port 8080
	fmt.Println("Starting HTTP server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error in Starting Server", err)
	}

}
