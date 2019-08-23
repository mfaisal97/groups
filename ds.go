package groups

import "golang.org/x/crypto/ed25519"

//alternatives
type Role string

// type RequestType string
//type MembershipRequestType string

//type Hash string
//type EncryptedHash string
type Data string

//type PublicKey []byte  //---> already implemented in crypto/ed25519
//type PrivateKey []byte //---> already implemented in crypto/ed25519

type DataBaseRequest struct {
	RequestNum int32
	Signature  Signature
}

type DataBaseResponse struct {
	Accepted bool
}

type DataBaseMessage struct {
	MessageStatus MessageStatus
}

type GroupCreationRequest struct {
	DataBaseRequest
	Group BasicGroup
}

type GroupCreationResponse struct {
	Accepted bool
}

type GroupCreationMessage struct {
	DataBaseMessage
	Request  GroupCreationRequest
	Response GroupCreationResponse
}

type MembershipRequest struct {
	DataBaseRequest
	UserID         string
	AffectedMember string
	AffectedRole   Role
	MembershipRequestType
	GroupName string
}

type MembershipResponse struct {
	RequestNumber int32
	RequestHash   [32]byte
	UserID        string
	Accepted      bool
	Signature     Signature
}

type MembershipMessage struct {
	DataBaseMessage
	Request  MembershipRequest
	Response MembershipResponse
}

type UserRequest struct {
	DataBaseRequest
	UserID     string
	PubllicKey ed25519.PublicKey
}

type UserResponse struct {
	Accepted bool
}

type UserMessage struct {
	DataBaseMessage
	Request  UserRequest
	Response UserResponse
}

// type Membership struct {
// 	types map[string]Role
// }

// type Role struct {
// 	members []UserID
// }

type Signature struct {
	Hash          [32]byte
	Encryptedhash []byte
}

//enums
type MessageStatus int

const (
	Received           MessageStatus = 0
	Confirmed          MessageStatus = 1
	Sent               MessageStatus = 2
	Succeeded          MessageStatus = 3
	ConfirmationFailed MessageStatus = 4
	Failed             MessageStatus = 5
)

type MembershipRequestType int

const (
	Join   MembershipRequestType = 0
	Remove MembershipRequestType = 1
)
