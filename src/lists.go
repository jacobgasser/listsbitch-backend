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
	list.Author = GetIDFromUsername(author)
	list.CreatedAt = time.Now()
	list.UpdatedAt = time.Now()
	if !DB.Where("title = ? AND author = ?", listName, list.Author).First(&List{}).RecordNotFound() {
		return fmt.Errorf("List with name %s already exists in %s's account", listName, author)
	}
	DB.Create(list)
	return nil
}