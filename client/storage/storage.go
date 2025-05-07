package storage

import (
	"bytes"
	"client/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var SecurityCheck = true

func checkSecurity(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("error parsing url: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("must use http or https scheme, current is %s", u.Scheme)
	}

	if SecurityCheck && u.Scheme != "https" {
		return fmt.Errorf("must use https scheme, current is %s", u.Scheme)
	}

	return nil
}

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
	url := common.GetBackendUrl()
	err := checkSecurity(url)
	if err != nil {
		return err
	}

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

	req, err := http.NewRequest(http.MethodPost, url+"/items", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(common.GetConf().Username, common.GetConf().Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error saving item, response code: %d", resp.StatusCode)
	}
	return nil
}

func GetItems() ([]common.ItemWithID, error) {
	url := common.GetBackendUrl()
	err := checkSecurity(url)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, url+"/items", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(common.GetConf().Username, common.GetConf().Password)

	resp, err := (&http.Client{}).Do(req)
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

	url := common.GetBackendUrl()
	err := checkSecurity(url)
	if err != nil {
		return i, err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/items/%d", url, id), nil)
	if err != nil {
		return i, err
	}
	req.SetBasicAuth(common.GetConf().Username, common.GetConf().Password)

	resp, err := (&http.Client{}).Do(req)
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
