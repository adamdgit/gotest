package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (app *app) home(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Print("Method: ", r.Method, "\n----------------------------\n",
		"Body: ", r.Body, "\n----------------------------\n", "Header: ", r.Header)

	// t, err := template.ParseFiles("../frontend/index.html")
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// t.Execute(w, map[string]any{"posts": posts})

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "Error encoding JSON data", 500)
	}
}
