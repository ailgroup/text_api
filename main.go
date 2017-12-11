package main

import (
	"log"
	"net/http"

	"github.com/jbowles/text_api/handlers"
)

var (
	port = ":8080"
)

//http://localhost:8080/detect_browse?text=some sentence to detect
//http://localhost:8080/detect?text=some sentence to detect
func main() {
	log.Printf("running server on port: %s", port)
	http.HandleFunc("/detect_browse", handlers.DetectBrowse)
	http.HandleFunc("/detect", handlers.Detect)
	http.HandleFunc("/sms_spam_predict", handlers.SmsSpamClassify)
	http.Handle("/favicon", http.NotFoundHandler())
	http.ListenAndServe(port, nil)
}
