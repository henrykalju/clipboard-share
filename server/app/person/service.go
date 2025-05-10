package person

import (
	"clipboard-share-server/db"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func GetPersonIDFromRequest(w http.ResponseWriter, r *http.Request) (int32, error) {
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("basic auth missing or empty")
	}

	person, err := db.Q.GetPerson(r.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusUnauthorized)
			return 0, errors.New("person does not exist")
		} else {
			fmt.Printf("Error getting person: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return 0, errors.New("error getting person")
		}
	}

	if person.Password != password {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("wrong password")
	}

	return person.ID, nil
}
