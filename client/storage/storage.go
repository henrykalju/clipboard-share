package storage

import (
	"bytes"
	"client/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BACKEND_URL = "http://localhost:8080/items"

type storageItem struct {
	ID      int32
	Type    string
	Content string
	Data    []data
}

type data struct {
	Format string
	Data   []byte
}

func SaveItem(i *common.Item) error {
	si := storageItem{
		Type:    i.Type.Text,
		Content: i.Text,
		Data:    []data{},
	}
	for _, v := range i.Values {
		si.Data = append(si.Data, data{
			Format: v.Format,
			Data:   v.Data,
		})
	}

	body, err := json.Marshal(&si)
	if err != nil {
		return err
	}

	resp, err := http.Post(BACKEND_URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error saving item, response code: %d", resp.StatusCode)
	}
	return nil
}

func GetItems() ([]common.ItemWithID, error) {
	resp, err := http.Get(BACKEND_URL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting items from storage, resp code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sitems []storageItem
	err = json.Unmarshal(body, &sitems)
	if err != nil {
		return nil, err
	}

	items := make([]common.ItemWithID, 0)
	for _, v := range sitems {
		i := common.ItemWithID{
			ID: v.ID,
			//Type: common.GetType(v.Type),
			Item: common.Item{
				Type:   common.GetType(v.Type),
				Text:   v.Content,
				Values: []common.Value{},
			},
		}

		for _, v2 := range v.Data {
			i.Values = append(i.Values, common.Value{
				Format: v2.Format,
				Data:   v2.Data,
			})
		}

		items = append(items, i)
	}

	return items, nil
}

func GetItemByID(id int32) (common.Item, error) {
	var i common.Item

	resp, err := http.Get(fmt.Sprintf("%s/%d", BACKEND_URL, id))
	if err != nil {
		return i, err
	}

	if resp.StatusCode != http.StatusOK {
		return i, fmt.Errorf("error getting items from storage, resp code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return i, err
	}

	var si storageItem
	err = json.Unmarshal(body, &si)
	if err != nil {
		return i, err
	}

	i = common.Item{
		Type:   common.GetType(si.Type),
		Text:   si.Content,
		Values: []common.Value{},
	}

	for _, v2 := range si.Data {
		i.Values = append(i.Values, common.Value{
			Format: v2.Format,
			Data:   v2.Data,
		})
	}

	return i, nil
}
