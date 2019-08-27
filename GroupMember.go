package groups

import (
	"math/rand"
	"time"
)

type User struct {
	UserInfo
	GroupMemberInfos map[string]GroupMemberInfo
}

//fetch requests in groups
func (user *User) UpdatePendingRequests(groupName string, requests []Message) bool {

	if info, exist := user.GroupMemberInfos[groupName]; exist {
		info.PendingRequests = requests
		user.GroupMemberInfos[groupName] = info
		return true
	}

	return false
}

//Add Group Info
func (user *User) AddGroupMemberInfo(groupName string, signRequest func(request Request, userInfo UserInfo) interface{}, signResponse func(response Response, userInfo UserInfo, userIDs []string) interface{}, requests []Message) {
	if user.GroupMemberInfos == nil {
		user.GroupMemberInfos = make(map[string]GroupMemberInfo)
	}
	var info GroupMemberInfo
	info.SignRequest = signRequest
	info.SignResponse = signResponse
	info.PendingRequests = requests
	user.GroupMemberInfos[groupName] = info
}

// create response to a fetched request
func (user *User) CreateResponseMessage(groupName string, requestNumber int, answer bool) ResponseMessage {
	if info, exist := user.GroupMemberInfos[groupName]; exist {
		if requestNumber >= 0 && requestNumber < len(info.PendingRequests) {
			meessage := info.PendingRequests[requestNumber]

			var responseMessage ResponseMessage
			responseMessage.Answer = answer
			responseMessage.requestHash = GenerateHash(meessage.RequestMessage)
			responseMessage.Signature = info.SignResponse(responseMessage.Response, user.UserInfo, meessage.AuthorizedUsers)
			return responseMessage
		}

	}
	return ResponseMessage{}
}

//create request for a group
func (user *User) CreateRequest(groupName string, requestName string, args ...interface{}) RequestMessage {
	if info, exist := user.GroupMemberInfos[groupName]; exist {
		var requestMessage RequestMessage
		requestMessage.RequestName = requestName
		requestMessage.RequestArgs = args
		requestMessage.RandomIdentifier = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000)
		requestMessage.Signature = info.SignRequest(requestMessage.Request, user.UserInfo)
	}
	return RequestMessage{}
}

//set sign request
func (user *User) SetSignRequestFunc(groupName string, signRequest func(request Request, userInfo UserInfo) interface{}) bool {

	if info, exist := user.GroupMemberInfos[groupName]; exist {
		info.SignRequest = signRequest
		user.GroupMemberInfos[groupName] = info
		return true
	}

	return false
}

//set sign request
func (user *User) SetSignResponseFunc(groupName string, signResponse func(response Response, userInfo UserInfo, userIDs []string) interface{}) bool {

	if info, exist := user.GroupMemberInfos[groupName]; exist {
		info.SignResponse = signResponse
		user.GroupMemberInfos[groupName] = info
		return true
	}

	return false
}

/*
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
*/
