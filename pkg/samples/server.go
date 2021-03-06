package samples

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"os"
)

var mu sync.Mutex
var count int

func Server() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/stop", stopper)
	log.Printf("Starting the goJam Server!")
	log.Fatal(http.ListenAndServe(":7777", nil))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()

	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}

	response := []byte("Hello World! - How are you ?")
	w.Write(response)
	return
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
	return
}

// stopper stops the running server.
func stopper(w http.ResponseWriter, r *http.Request) {
	log.Printf("Stopping the goJam Server!")
	stop := make(chan os.Signal, 1)
	<-stop
}
