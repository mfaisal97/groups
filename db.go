package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"crypto/sha256"

	"golang.org/x/crypto/ed25519"
)

//stores all valid reqeusts and replies and verify sender.
type DataBase struct {
	FileName string

	Groups []Group
	Users  map[UserID]ed25519.PublicKey

	GroupCreationMessages []GroupCreationMessage
	GroupCreationIndex    int32

	MembershipsMessages []MembershipMessage
	MembershipIndex     int32

	UserMessages []UserMessage
	UsersIndex   int32
}

func InitializeDataBase() DataBase {
	var db DataBase

	db.Users = make(map[UserID]ed25519.PublicKey, 0)
	db.GroupCreationIndex = 0
	db.MembershipIndex = 0
	db.UsersIndex = 0

	return db
}

func (db *DataBase) ResetDataBase() {
	fileName := db.FileName
	*db = InitializeDataBase()
	db.SaveData(fileName)
}

func (db *DataBase) LoadData(fileName string) {
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

func (db *DataBase) SaveData(fileName string) {
	db.FileName = fileName
	file, _ := json.MarshalIndent(db, "", " ")
	_ = ioutil.WriteFile(fileName, file, 0644)
	return
}

func (db *DataBase) GetRequestNumber(userRequest UserRequest) int32 {
	userRequest.DataBaseRequest.RequestNum = db.UsersIndex

	newUserMessage := UserMessage{}

	newUserMessage.Request = userRequest

	jsonString, err := json.Marshal(newUserMessage.Request)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return -1
	}
	newUserMessage.Request.Signature.Hash = sha256.Sum256([]byte(jsonString))

	newUserMessage.MessageStatus = 0

	db.UserMessages = append(db.UserMessages, newUserMessage)

	db.IncreaseUsersIndex()

	defer db.SaveData(db.FileName)
	return newUserMessage.Request.RequestNum
}

func (db *DataBase) IncreaseUsersIndex() {
	db.UsersIndex = db.UsersIndex + 1
}

func (db *DataBase) ConfirmRequest(userRequest UserRequest) UserResponse {
	requestNumber := userRequest.RequestNum

	userMessage := db.UserMessages[requestNumber]

	if userMessage.Request.Signature.Hash == userRequest.Signature.Hash &&
		ed25519.Verify(userMessage.Request.PubllicKey, []byte(userMessage.Request.Signature.Hash[:]), userRequest.Signature.Encryptedhash) {
		db.UserMessages[requestNumber].Request.Signature.Encryptedhash = userRequest.Signature.Encryptedhash
		db.UserMessages[requestNumber].Response = UserResponse{true}
		db.UserMessages[requestNumber].DataBaseMessage.MessageStatus = Succeeded
		defer db.SaveData(db.FileName)

		return db.UserMessages[requestNumber].Response
	}

	return UserResponse{false}
}
