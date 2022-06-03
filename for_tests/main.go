package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func echoAwake(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "I'm awake ")
}

func echoSleepy(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Second * 5)
	fmt.Fprintf(w, "I'm kinda sleepy")
}

func echoGreetings(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hey there!")
}

func main() {
	http.HandleFunc("/", echoGreetings)
	http.HandleFunc("/awake", echoAwake)
	http.HandleFunc("/sleepy", echoSleepy)
	log.Fatal(http.ListenAndServe(":4321", nil))
}
