package main

import (
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm"
	"time"
)

func CreateList(author string, listName string) error {
	list := List{}
	list.Title = listName
	list.ID = uuid.New().String()
	list.Author = author
	list.CreatedAt = time.Now()
	list.UpdatedAt = time.Now()
	if !DB.Where("title = ? AND author = ?", listName, list.Author).First(&List{}).RecordNotFound() {
		return fmt.Errorf("lists: list with name %s already exists in %s's account", listName, author)
	}
	DB.Create(list)
	return nil
}

func CreateListItem(author string, listID string, listItemContent string) error {
	listItem := ListItem{}
	listItem.ID = uuid.New().String()
	listItem.ListID = listID
	listItem.Content = listItemContent
	DB.Create(&listItem)
	return nil
}

func GetListItems(listID string) ([]ListItem , error) {
	listItems := []ListItem{}
	if DB.Find(&listItems, "list_id = ?", listID).RecordNotFound() {
		return listItems, fmt.Errorf("lists: no list with ID of %s exists", listID)
	}
	return listItems, nil
}