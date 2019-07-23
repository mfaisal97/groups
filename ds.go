package main

import "golang.org/x/crypto/ed25519"

//alternatives
type Role string
type RequestType string
type MembershipRequestType string

//type Hash string
//type EncryptedHash string
type Data string

//type PublicKey []byte  //---> already implemented in crypto/ed25519
//type PrivateKey []byte //---> already implemented in crypto/ed25519

type Group struct {
	id          string
	name        string
	description string
	data        map[RequestType]Data

	members               []UserID
	memberships           map[Role][]UserID
	authorizations        map[RequestType][]Role
	membersauthorizations map[UserID][]RequestType
}

type DataBaseRequest struct {
	rquestNum int32
	signature Signature
}

type DataBaseMessage struct {
	messageStatus MessageStatus
}

type GroupCreationRequest struct {
	DataBaseRequest
	group Group
}

type GroupCreationResponse struct {
	num       int32
	signature Signature
}

type GroupCreationMessage struct {
	DataBaseMessage
	request  GroupCreationRequest
	response GroupCreationResponse
}

type MembershipRequest struct {
	DataBaseRequest
}

type MembershipResponse struct {
}

type MembershipMessage struct {
	DataBaseMessage
	request  MembershipRequest
	response MembershipResponse
}

type UserRequest struct {
	DataBaseRequest
	name       UserID
	publlicKey ed25519.PublicKey
}

type UserResponse struct {
	Accepted bool
}

type UserMessage struct {
	DataBaseMessage
	request  *UserRequest
	response *UserResponse
}

// type Membership struct {
// 	types map[string]Role
// }

// type Role struct {
// 	members []UserID
// }

type Signature struct {
	hash          [32]byte
	encryptedhash []byte
}

//enums
type MessageStatus int

const (
	Received  MessageStatus = 0
	Sent      MessageStatus = 1
	Succeeded MessageStatus = 2
	Failed    MessageStatus = 3
)
