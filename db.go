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

	Groups map[string]Group
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

	if _, exist := db.Users[userRequest.UserID]; exist {
		fmt.Printf("This UserID Already Exists: %s", userRequest.UserID)
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
		db.Users[userRequest.UserID] = userRequest.PubllicKey
		defer db.SaveData(db.FileName)

		return db.UserMessages[requestNumber].Response
	}

	return UserResponse{false}
}

func (db *DataBase) GetGroupCreationRequestNumber(groupCreationRequest GroupCreationRequest) int32 {
	groupCreationRequest.DataBaseRequest.RequestNum = db.GroupCreationIndex

	newGroupCreationMessage := GroupCreationMessage{}
	newGroupCreationMessage.Request = groupCreationRequest

	jsonString, err := json.Marshal(newGroupCreationMessage.Request)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return -1
	}
	newGroupCreationMessage.Request.Signature.Hash = sha256.Sum256([]byte(jsonString))

	newGroupCreationMessage.MessageStatus = 0

	db.GroupCreationMessages = append(db.GroupCreationMessages, newGroupCreationMessage)

	db.IncreaseGroupCreationIndex()

	defer db.SaveData(db.FileName)
	return newGroupCreationMessage.Request.RequestNum
}

func (db *DataBase) IncreaseGroupCreationIndex() {
	db.GroupCreationIndex = db.GroupCreationIndex + 1
}

func (db *DataBase) ConfirmGroupCreationRequest(groupCreationRequest GroupCreationRequest) GroupCreationResponse {

	requestNumber := groupCreationRequest.RequestNum

	groupCreationMessage := db.GroupCreationMessages[requestNumber]

	if groupCreationMessage.Request.Signature.Hash == groupCreationRequest.Signature.Hash &&
		ed25519.Verify(db.Users[groupCreationMessage.Request.Group.Creator], []byte(groupCreationMessage.Request.Signature.Hash[:]), groupCreationRequest.Signature.Encryptedhash) {

		if _, exist := db.Groups[groupCreationRequest.Group.Name]; exist {
			fmt.Printf("This Group Already Exists: %s", groupCreationRequest.Group.Name)
			db.GroupCreationMessages[requestNumber].Response = GroupCreationResponse{false}
			db.GroupCreationMessages[requestNumber].MessageStatus = Failed
			return GroupCreationResponse{false}
		}

		//ensuring that each user exits once
		set := make(map[UserID]bool)
		for k := range groupCreationRequest.Group.Members {
			set[groupCreationRequest.Group.Members[k]] = true
		}
		if len(set) != len(groupCreationRequest.Group.Members) {
			fmt.Println("Error: Request Contains Reptited Members")
			db.GroupCreationMessages[requestNumber].Response = GroupCreationResponse{false}
			db.GroupCreationMessages[requestNumber].MessageStatus = Failed
			return GroupCreationResponse{false}
		}

		//ensuring all memebers are users
		for k := range groupCreationRequest.Group.Members {
			if _, exist := db.Users[groupCreationRequest.Group.Members[k]]; !exist {
				fmt.Println("Error: Request Contains not registered Members")
				db.GroupCreationMessages[requestNumber].Response = GroupCreationResponse{false}
				db.GroupCreationMessages[requestNumber].MessageStatus = Failed
				return GroupCreationResponse{false}
			}
		}

		//ensuring that the members in each role are in the group members

		for key, val := range groupCreationRequest.Group.Memberships {

			set2 := make(map[UserID]bool)
			for k := range val {
				set[val[k]] = true
			}
			if len(set2) != len(val) {
				fmt.Println("Error: Request Contains Reptited Members in Role:\t", key)
				db.GroupCreationMessages[requestNumber].Response = GroupCreationResponse{false}
				db.GroupCreationMessages[requestNumber].MessageStatus = Failed
				return GroupCreationResponse{false}
			}

			//ensuring all role memebers are members
			for k := range val {
				if exist := set[val[k]]; !exist {
					fmt.Println("Error: Role members Contains non Members in Role:\t", key)
					db.GroupCreationMessages[requestNumber].Response = GroupCreationResponse{false}
					db.GroupCreationMessages[requestNumber].MessageStatus = Failed
					return GroupCreationResponse{false}
				}
			}

		}

		db.GroupCreationMessages[requestNumber].Request.Signature.Encryptedhash = groupCreationRequest.Signature.Encryptedhash
		db.GroupCreationMessages[requestNumber].Response = GroupCreationResponse{true}
		db.GroupCreationMessages[requestNumber].DataBaseMessage.MessageStatus = Succeeded
		db.Groups[groupCreationRequest.Group.Name] = groupCreationRequest.Group
		defer db.SaveData(db.FileName)

		return db.GroupCreationMessages[requestNumber].Response
	}

	return GroupCreationResponse{false}
}
