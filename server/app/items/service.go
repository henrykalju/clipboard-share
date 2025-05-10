package items

import (
	"clipboard-share-server/db"
	"context"
	"errors"
	"fmt"
)

type ItemWithData struct {
	db.Item
	Data []db.Datum
}

func getItemsWithDataByPerson(personID int32) ([]ItemWithData, error) {
	ctx := context.Background()
	r := make([]ItemWithData, 0)

	items, err := db.Q.GetItemsByPerson(ctx, personID)
	if err != nil {
		return nil, err
	}

	for _, v := range items {
		item := ItemWithData{Item: v}

		data, err := db.Q.GetDataByItem(ctx, item.ID)
		if err != nil {
			return nil, err
		}

		item.Data = append(item.Data, data...)

		r = append(r, item)
	}

	return r, nil
}

func getItemWithDataByIdAndPerson(itemID int32, personID int32) (ItemWithData, error) {
	ctx := context.Background()
	r := ItemWithData{}

	params := db.GetItemByIdAndPersonParams{
		ID:       itemID,
		PersonID: personID,
	}
	item, err := db.Q.GetItemByIdAndPerson(ctx, params)
	if err != nil {
		return r, err
	}

	r.Item = item
	data, err := db.Q.GetDataByItem(ctx, item.ID)
	if err != nil {
		return r, err
	}

	r.Data = data

	return r, nil
}

func (i ItemWithData) validateInsert() error {
	if len(i.Content) == 0 {
		return errors.New("item must have content string")
	}
	return nil
}

func insertItem(item ItemWithData) (ItemWithData, error) {
	err := item.validateInsert()
	if err != nil {
		return item, err
	}

	ctx := context.Background()
	tx, err := db.NewTx(ctx)
	if err != nil {
		return item, err
	}
	q := db.Q.WithTx(tx)
	defer func() {
		if err != nil {
			fmt.Println("Rolling back item addition")
			err = tx.Rollback(ctx)
			if err != nil {
				fmt.Printf("Error rolling back adding item: %s\n", err.Error())
			}
		}
	}()

	itemParams := db.InsertItemParams{
		PersonID: item.PersonID,
		Content:  item.Content,
		Type:     item.Type,
	}

	i, err := q.InsertItem(ctx, itemParams)
	if err != nil {
		return item, err
	}

	returnedItem := ItemWithData{Item: i}

	for _, v := range item.Data {
		dataParams := db.InsertDataParams{
			ItemID: i.ID,
			Format: v.Format,
			Data:   v.Data,
		}
		d, err := q.InsertData(ctx, dataParams)
		if err != nil {
			return item, err
		}

		returnedItem.Data = append(returnedItem.Data, d)
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("Error committing adding item: %s\n", err.Error())
		return item, err
	}

	go checkSizes(returnedItem.PersonID)

	return returnedItem, nil
}

func checkSizes(personID int32) {
	err := db.Q.CheckSizes(context.Background(), db.CheckSizesParams{
		PersonID:  personID,
		Threshold: 1024 * 1024,
	})
	if err != nil {
		fmt.Printf("Error deleting items that cross threshold of 1MB: %s\n", err.Error())
	}
}
