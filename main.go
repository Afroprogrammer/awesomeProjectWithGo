package main

import (
	"awesomeProject/news"
	"fmt"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, nil)
	if err != nil {
		return
	}
}

// method below used for search
func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}
		fmt.Println("Search Query is: ", searchQuery)
		fmt.Println("Page is: ", page)
	}
}
func main() {
	//reading port from env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	apikey := os.Getenv("NEWS_API_KEY")
	if apikey == "" {
		log.Fatal("Env: apikey must be set")
	}
	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apikey, 20)

	fs := http.FileServer(http.Dir("assets"))
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/search", searchHandler(newsapi))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
