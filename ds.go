package main

//alternatives
type Role string
type RequestType string
type MembershipRequestType string
type Hash string
type EncryptedHash string
type Data string
type PublicKey string
type PrivateKey string

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

type GroupCreationRequest struct {
	num       int32
	group     Group
	signature Signature
}

type GroupCreationResponse struct {
	num       int32
	signature Signature
}

type GroupCreationMessage struct {
	request  GroupCreationRequest
	response GroupCreationResponse
}

type MembershipRequest struct {
}

type MembershipResponse struct {
}

type MembershipMessage struct {
	request  MembershipRequest
	response MembershipResponse
}

// type Membership struct {
// 	types map[string]Role
// }

// type Role struct {
// 	members []UserID
// }

type Signature struct {
	requestNum    int
	hash          Hash
	encryptedhash EncryptedHash
}

//enums
type MessageStatus int

const (
	Received  MessageStatus = 0
	Sent      MessageStatus = 1
	Succeeded MessageStatus = 2
	Failed    MessageStatus = 3
)
