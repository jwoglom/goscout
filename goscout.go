package main

import (
	"./app"
)

func main() {
	s := app.NewServer()

	s.Run()
	/*
		r := mux.NewRouter()
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "dummy webpage\n")
		})

		fmt.Println("running on :3000")
		http.ListenAndServe(":3000", r)
	*/
}
