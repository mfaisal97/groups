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
	UserID
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
	user.UserID = userid

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
	userFile, err := os.Open("Users/" + string(userid) + ".json")
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
	_ = ioutil.WriteFile("Users/"+string(userid)+".json", file, 0644)
	return
}

func (user *User) CreateGroup(group Group) {

	groupCreationRequest := GroupCreationRequest{}
	groupCreationRequest.Group = group

	requestNum := data.GetGroupCreationRequestNumber(groupCreationRequest)
	if requestNum == -1 {
		_ = fmt.Errorf("Couldnot Get Request Number for New Group Creation")
		return
	}
	fmt.Println("Got Request Number:\t\t", requestNum)

	groupCreationRequest.RequestNum = requestNum

	jsonString, err := json.Marshal(groupCreationRequest)
	if err != nil {
		_ = fmt.Errorf("Error: %s", err)
		return
	}
	groupCreationRequest.Signature.Hash = sha256.Sum256([]byte(jsonString))
	groupCreationRequest.Signature.Encryptedhash = ed25519.Sign(user.PrivateKey, []byte(groupCreationRequest.Signature.Hash[:]))

	if !data.ConfirmGroupCreationRequest(groupCreationRequest).Accepted {
		_ = fmt.Errorf("Error: Signature was not accepted")
		return
	}

	fmt.Println("New Group was created !")
	return
}

func (user *User) SendMembershipRequest(membershipRequestType MembershipRequestType, group string, affectedMember UserID, affectedRole Role) {

	membershipRequest := MembershipRequest{}
	membershipRequest.AffectedMember = affectedMember
	membershipRequest.AffectedRole = affectedRole
	membershipRequest.UserID = user.UserID
	membershipRequest.GroupName = group
	membershipRequest.MembershipRequestType = membershipRequestType

	requestNum := data.GetMembershipRequestNumber(membershipRequest)
	if requestNum == -1 {
		_ = fmt.Errorf("Couldnot Get Request Number for New Group Creation")
		return
	}
	fmt.Println("Got Request Number:\t\t", requestNum)

	membershipRequest.RequestNum = requestNum

	jsonString, err := json.Marshal(membershipRequest)
	if err != nil {
		_ = fmt.Errorf("Error: %s", err)
		return
	}
	membershipRequest.Signature.Hash = sha256.Sum256([]byte(jsonString))
	membershipRequest.Signature.Encryptedhash = ed25519.Sign(user.PrivateKey, []byte(membershipRequest.Signature.Hash[:]))

	if data.ConfirmMembershipRequest(membershipRequest) != Confirmed {
		_ = fmt.Errorf("Error: Signature was not accepted")
		return
	}

	fmt.Println("New Membership Request was created !")
	return
}

func (user *User) GetPendingRequests(groupName string) ([]int32, []MembershipRequest) {
	return data.GetPendingRequests(user.UserID, groupName)
}

func (user *User) SendMembershipReesponse(requestNumber int32, accepted bool) {

	membershipRequest := data.GetMembershipRequest(requestNumber)

	membershipResponse := MembershipResponse{}
	membershipResponse.RequestNumber = requestNumber
	membershipResponse.RequestHash = membershipRequest.Signature.Hash
	membershipResponse.UserID = user.UserID
	membershipResponse.Accepted = accepted

	jsonString, err := json.Marshal(membershipResponse)
	if err != nil {
		_ = fmt.Errorf("Error: %s", err)
		return
	}
	membershipResponse.Signature.Hash = sha256.Sum256([]byte(jsonString))
	membershipResponse.Signature.Encryptedhash = ed25519.Sign(user.PrivateKey, []byte(membershipResponse.Signature.Hash[:]))

	messageStatus := data.ConfirmMembershipResponse(membershipResponse)
	if messageStatus == ConfirmationFailed {
		_ = fmt.Errorf("Error: Signature was not accepted")
		return
	}

	fmt.Println("Membership Request Got New Status !")
	return
}
