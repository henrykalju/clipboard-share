package items

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	/*personIdS := r.Header.Get("person_id")
	personId64, err := strconv.ParseInt(personIdS, 10, 32)
	if err != nil {
		fmt.Printf("Error parsing person_id: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/
	personId64 := 1

	items, err := getItemsWithDataByPerson(int32(personId64))
	if err != nil {
		fmt.Printf("Error getting items: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&items)
	if err != nil {
		fmt.Printf("Error marshalling items into response: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err.Error())
	}
}

func GetItemByID(w http.ResponseWriter, r *http.Request) {
	/*personIdS := r.Header.Get("person_id")
	personId64, err := strconv.ParseInt(personIdS, 10, 32)
	if err != nil {
		fmt.Printf("Error parsing person_id: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/
	personId64 := 1

	idString := r.PathValue("id")
	id64, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := getItemWithDataByIdAndPerson(int32(id64), int32(personId64))
	if err != nil {
		fmt.Printf("Error getting item with data by id and person: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&i)
	if err != nil {
		fmt.Printf("Error marshalling item with data into response: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err.Error())
	}
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var item ItemWithData

	err = json.Unmarshal(body, &item)
	if err != nil {
		fmt.Printf("Error unmarshalling body: %s\n, %s\n", err.Error(), string(body))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	item.PersonID = 1

	item, err = insertItem(item)
	if err != nil {
		// TODO if errors.As Validation error, 400, else 500
		fmt.Printf("Error inserting item: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&item)
	if err != nil {
		fmt.Printf("Error marshalling response: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
