package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("view"))
	http.Handle("/", fs)
	fmt.Println("ready")
	http.ListenAndServe(":8000", nil)
}
