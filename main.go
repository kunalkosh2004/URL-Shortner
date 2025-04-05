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

type URL struct{
	ID string `json:"id"`
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func createURL(originalURL string) string{
	shortUrl := generateShortURL(originalURL)
	id := shortUrl
	urlDB[id] = URL{
		ID: id,
        ShortURL: shortUrl,
        OriginalURL: originalURL,
        CreationDate: time.Now(),
	}
	return shortUrl
}

func getURL(id string)(URL, error){
	url,ok := urlDB[id]
	if !ok{
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"hello World!")
}

func shortURLHandler(w http.ResponseWriter,r *http.Request){
	var data struct{
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	shortURL_ := createURL(data.URL)
	// fmt.Fprintf(w,shortURL)
	response := struct{
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request){
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main()  {
	fmt.Println("URL-Shortener")
	// url := "https://github.com/kunalkosh2004"
	// asd := generateShortURL(url)
	// fmt.Println("short: ",asd)

	// Register the handler function to handle all requests to the root URL
	http.HandleFunc("/",handler)
	http.HandleFunc("/shorten", shortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	// Start the http server
	fmt.Println("Server starting on 8080 port")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
        fmt.Println("Error on starting http server")
    }
}