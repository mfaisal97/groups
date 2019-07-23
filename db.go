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

	db.Users = make(map[UserID]ed25519.PublicKey)
	db.GroupCreationIndex = 0
	db.MembershipIndex = 0
	db.UsersIndex = 0

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

func (db DataBase) GetRequestNumber(userRequest UserRequest) int32 {
	userRequest.DataBaseRequest.rquestNum = db.UsersIndex

	newUserMessage := UserMessage{}

	newUserMessage.request = userRequest

	jsonString, err := json.Marshal(newUserMessage.request)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return -1
	}
	newUserMessage.request.signature.hash = sha256.Sum256([]byte(jsonString))

	newUserMessage.messageStatus = 0
	fmt.Println(db)
	fmt.Println("Adding:", newUserMessage)
	db.UserMessages = append(db.UserMessages, newUserMessage)
	fmt.Println("After:", db.UserMessages)
	defer db.IncreaseUsersIndex()
	return newUserMessage.request.rquestNum
}

func (db DataBase) IncreaseUsersIndex() {
	db.UsersIndex = db.UsersIndex + 1
	fmt.Println("After After:", db.UserMessages)
}

func (db DataBase) ConfirmRequest(userRequest UserRequest) UserResponse {
	requestNumber := userRequest.rquestNum

	fmt.Println("Using RN:", requestNumber)
	fmt.Println("In:", db.UserMessages)
	userMessage := db.UserMessages[requestNumber]

	if userMessage.request.signature.hash == userRequest.signature.hash &&
		ed25519.Verify(userMessage.request.publlicKey, []byte(userMessage.request.signature.hash[:]), userRequest.signature.encryptedhash) {
		db.UserMessages[requestNumber].request.signature.encryptedhash = userRequest.signature.encryptedhash
		db.UserMessages[requestNumber].response = UserResponse{true}
		db.UserMessages[requestNumber].DataBaseMessage.messageStatus = Succeeded

		return db.UserMessages[requestNumber].response
	}

	return UserResponse{false}
}
