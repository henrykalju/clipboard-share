package storage

import (
	"bytes"
	"client/common"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	SecurityCheck = true
	ws            *websocket.Conn
)

func checkSecurity(raw string) (string, error) {
	if !strings.Contains(raw, "://") {
		raw = "x://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("error parsing url: %w", err)
	}

	if u.Scheme == "x" {
		u.Scheme = ""
	}

	if u.Scheme == "" {
		u.Scheme = "http"
		if SecurityCheck {
			u.Scheme += "s"
		}
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("must use http or https scheme, current is %s", u.Scheme)
	}

	if SecurityCheck && u.Scheme != "https" {
		return "", fmt.Errorf("must use https scheme, current is %s", u.Scheme)
	}

	return u.String(), nil
}

func RestartWebsocket(c chan bool) error {
	if ws != nil {
		err := ws.Close()
		if err != nil {
			return err
		}
		ws = nil
	}

	return StartWebsocket(c)
}

func StartWebsocket(c chan bool) error {
	if ws != nil {
		panic("Websocket is not nil and trying to start one")
	}

	conf := common.GetConf()
	raw := conf.BackendUrl
	if !strings.Contains(raw, "://") {
		raw = "x://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme == "x" {
		u.Scheme = ""
	}

	if u.Scheme == "" {
		u.Scheme = "ws"
		if SecurityCheck {
			u.Scheme += "s"
		}
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		return fmt.Errorf("must use ws or wss scheme, current is %s", u.Scheme)
	}

	if SecurityCheck && u.Scheme != "wss" {
		return fmt.Errorf("must use wss scheme, current is %s", u.Scheme)
	}

	header := http.Header{}
	header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(conf.Username+":"+conf.Password)))

	conn, resp, err := websocket.DefaultDialer.Dial(u.String()+"/ws", header)
	if err != nil {
		if resp == nil {
			return fmt.Errorf("error dialing ws: %w", err)
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("incorrect username or password")
		}
	}
	if conn == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	go func() {
		for {
			msgType, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading websocket: %s\n", err.Error())
				c <- false
				return
			}
			if msgType != websocket.TextMessage {
				fmt.Printf("Websocket sent not textmessage: %d\n", msgType)
				continue
			}
			if string(message) == "HISTORY" {
				c <- true
			}
		}
	}()

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
	url, err := checkSecurity(common.GetBackendUrl())
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
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("incorrect username or password")
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("response code: %d", resp.StatusCode)
	}
	return nil
}

func GetItems() ([]common.ItemWithID, error) {
	url, err := checkSecurity(common.GetBackendUrl())
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

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("incorrect username or password")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resp code: %d", resp.StatusCode)
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

	url, err := checkSecurity(common.GetBackendUrl())
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

	if resp.StatusCode == http.StatusUnauthorized {
		return i, fmt.Errorf("incorrect username or password")
	}

	if resp.StatusCode != http.StatusOK {
		return i, fmt.Errorf("resp code: %d", resp.StatusCode)
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
