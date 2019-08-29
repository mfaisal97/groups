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
