package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//stores all valid reqeusts and replies and verify sender.
type DataBase struct {
	FileName string
	Groups   []Group

	GroupCreationMessages []GroupCreationMessage
	MembershipsMessages   []MembershipMessage

	Users map[UserID]PublicKey
}

func InitializeDataBase() DataBase {
	var db DataBase
	db.Users = make(map[UserID]PublicKey)
	return db
}

func (db DataBase) ResetDataBase() {
	fileName := db.FileName
	db = InitializeDataBase()
	db.SaveData(fileName)
}

func (db DataBase) LoadData(fileName string) {
	db.FileName = fileName

	dbFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	byteValue, _ := ioutil.ReadAll(dbFile)
	json.Unmarshal(byteValue, &db)
	db.FileName = fileName
	defer dbFile.Close()
	return
}

func (db DataBase) SaveData(fileName string) {
	db.FileName = fileName
	file, _ := json.MarshalIndent(db, "", " ")
	_ = ioutil.WriteFile(fileName, file, 0644)
	return
}
