package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ed25519"
)

type UserID string

type User struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

func CreateNewUser(userid UserID) User {
	var user User
	var rand io.Reader
	pub, pv, err := ed25519.GenerateKey(rand)
	if err != nil {
		log.Fatal(err)
	}
	user.PublicKey = pub
	user.PrivateKey = pv

	userRequest := UserRequest{}
	userRequest.UserID = userid
	userRequest.PubllicKey = user.PublicKey

	requestNum := data.GetRequestNumber(userRequest)
	if requestNum == -1 {
		_ = fmt.Errorf("Couldnot Get Request Number for New User Creation")
		return user
	}
	fmt.Println("Got Request Number:\t\t", requestNum)

	userRequest.RequestNum = requestNum

	jsonString, err := json.Marshal(userRequest)
	if err != nil {
		_ = fmt.Errorf("Error: %s", err)
		return user
	}
	userRequest.Signature.Hash = sha256.Sum256([]byte(jsonString))
	userRequest.Signature.Encryptedhash = ed25519.Sign(pv, []byte(userRequest.Signature.Hash[:]))

	if !data.ConfirmRequest(userRequest).Accepted {
		_ = fmt.Errorf("Error: Signature was not accepted")
		return user
	}

	fmt.Println("New User was created !")
	user.SaveUser(userid)
	return user
}

func (user *User) LoadData(userid UserID) {
	userFile, err := os.Open(string(userid) + ".json")
	if err != nil {
		fmt.Println(err)
		return
	}

	byteValue, _ := ioutil.ReadAll(userFile)
	json.Unmarshal(byteValue, &user)
	defer userFile.Close()
	return
}

func (user *User) SaveUser(userid UserID) {
	file, _ := json.MarshalIndent(user, "", " ")
	_ = ioutil.WriteFile(string(userid)+".json", file, 0644)
	return
}
