package person

import (
	"clipboard-share-server/db"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		fmt.Printf("Basic auth missing or empty: %s, %s, %t\n", username, password, ok)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person, err := db.Q.GetPerson(r.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("person not exist")
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			fmt.Printf("Error getting person: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if person.Password != password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	type register struct {
		Username string
		Password string
	}

	var in register
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		fmt.Println("Error decoding register body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if in.Username == "" || in.Password == "" {
		fmt.Println("Username or password empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Q.GetPerson(r.Context(), in.Username)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("Error getting person: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Q.InsertPerson(r.Context(), db.InsertPersonParams{
		Username: in.Username,
		Password: in.Password,
	})
	if err != nil {
		fmt.Printf("Error inserting person: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
