package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("web3"))
	})
	http.ListenAndServe(":9003", nil)
}
