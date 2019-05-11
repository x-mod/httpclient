package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/error invoked")
		http.Error(w, "error ...", http.StatusBadRequest)
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/hello invoked")
		fmt.Fprintf(w, "hello world")
	})
	http.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sleep invoked")
		time.Sleep(500 * time.Millisecond)
		fmt.Fprintf(w, "awake ok")
	})
	log.Println("http serving at 127.0.0.1:12345 ...")
	for {
		if err := http.ListenAndServe("127.0.0.1:12345", nil); err != nil {
			log.Println("http serving failed: ", err)
		}
		log.Println("http serving at 127.0.0.1:12345 ...")
	}

}
