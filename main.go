package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"time"
)

var ApiUrl = "https://api.datamuse.com/words?md=s&"
var PORT = "8080"
var BasePath = "/api/"

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc(BasePath, apiHandler).Methods("GET", "OPTIONS").Queries("search", "")
	r.HandleFunc(BasePath, apiHandler).Methods("GET", "OPTIONS").Queries("search", "", "rel", "")
	r.HandleFunc("*", notFoundHandler)

	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	srv := &http.Server{
		Addr:         "0.0.0.0:" + PORT,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cor.Handler(r),
	}

	// Run server in a goroutine so that it doesn't block.
	go func() {
		fmt.Println("Starting server on port " + PORT)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

// Handler func for search endpoint
// Makes api call to Datamuse, sorts data, and returns response to user
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("OPTION", "*")
	// get search
	search := r.FormValue("search")

	// 400 if not found
	if len(search) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("search cannot be empty"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	searchType := "ml="
	// check for rel key
	if r.URL.Query().Has("rel") {
		searchType = "rel_syn="
	}

	// request to datamuse
	res, err := http.Get(ApiUrl + searchType + search)
	if err != nil {
		serverError(err, w)
		return
	}

	defer res.Body.Close()

	// read body
	b, err := io.ReadAll(res.Body)
	if err != nil {
		serverError(err, w)
		return
	}

	// unmarshal to slice of Words
	var words []Word
	err = json.Unmarshal(b, &words)
	if err != nil {
		serverError(err, w)
		return
	}

	// sort slice
	sort.Slice(words, func(p, q int) bool {
		// sort by numSyllables if the value isn't the same
		if words[p].NumSyllables != words[q].NumSyllables {
			return words[p].NumSyllables < words[q].NumSyllables
		}

		// else, sort by the word
		return words[p].Word < words[q].Word
	})

	// marshall slice into []byte
	resp, err := json.Marshal(words)
	if err != nil {
		serverError(err, w)
		return
	}

	// write response
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// Logs error and writes a 500 error
func serverError(err error, w http.ResponseWriter) {
	log.Println(err)

	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write([]byte("Internal Server Error"))
	if err != nil {
		log.Fatal(err)
	}
}

// Handler func for catch all route, just return 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Not found"))
	if err != nil {
		log.Fatal(err)
	}
}
