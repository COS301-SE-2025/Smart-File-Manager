package main

import "time"

type managedItem struct {
	itemID       string
	itemName     string
	itemPath     string // Use string for file paths; there is no built-in 'filepath' type
	itemTags     []tag
	locked       bool
	fileType     string
	creationDate time.Time
}

//composite functions
// func (m managedItem) createBackup() memento {

// }
// func (m managedItem) createBackup(savedBackup memento) {

// }

type Folder struct {
	managedItem
	containedItems []managedItem
}

//component functions
// func (f *Folder) addItem(newItem managedItem) {
// }

// func (f *Folder) removeItem(existingItem managedItem) bool {

// }

type File struct {
	managedItem
}

type tag struct {
	tagID   string
	tagName string
}
