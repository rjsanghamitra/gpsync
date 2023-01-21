package main

import (
	"fmt"
	"net/http"

	"github.com/rjsanghamitra/gpsync/handlers"
	"github.com/rjsanghamitra/gpsync/oauth"
)

func main() {
	oauth.Client()
	fmt.Println("Please open the link below in a browser:")
	fmt.Println("localhost:8080/")
	http.HandleFunc("/", handlers.Login)
	http.HandleFunc("/callback", handlers.OauthAndDownload)
	http.ListenAndServe(":8080", nil)
}
