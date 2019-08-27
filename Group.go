package groups

type Group struct {
	BasicGroup

	Messages        map[string]Message
	AllMessages     []interface{}
	PendingRequests map[string]map[string]EmptyStruct
}

func (group *Group) recordMessage(message interface{}) {
	group.AllMessages = append(group.AllMessages, message)
}

//Authorized members are specified on request delivery
//Execution is applied after authorized members approval and if the request type is stil valid

func (group *Group) AddRequest(requestMessage RequestMessage) RequestStatus {
	defer group.recordMessage(requestMessage)

	if group.IsRequestType(requestMessage.RequestName) {
		hash := GenerateHash(requestMessage.Request)
		if _, ok := group.Messages[string([]byte(hash[:]))]; !ok {
			verifyFunc, ok2 := group.VerifyRequest[requestMessage.RequestName]
			if (ok2 && verifyFunc(requestMessage.Request, requestMessage.Signature) == RequestConfirmed) || (!ok2 || group.DefaultVerifyRequest(requestMessage.Request, requestMessage.Signature) == RequestConfirmed) {
				{
					newMessage := Message{}
					newMessage.RequestMessage = requestMessage
					newMessage.Responses = make(map[string]ResponseMessage)
					newMessage.AuthorizedUsers = group.GetMembersForRequestType(requestMessage.RequestName)
					newMessage.MessageStatus = OnGoing
					group.Messages[string([]byte(hash[:]))] = newMessage
					if requests, exists := group.PendingRequests[requestMessage.RequestName]; exists {
						requests[string([]byte(hash[:]))] = EmptyStruct{}
						group.PendingRequests[requestMessage.RequestName] = requests
					} else {
						requests := make(map[string]EmptyStruct)
						requests[string([]byte(hash[:]))] = EmptyStruct{}
						group.PendingRequests[requestMessage.RequestName] = requests
					}
					return RequestReceived
				}
			}
		}
	}
	return RequestFailed
}

func (group *Group) AddResponse(responseMessage ResponseMessage) MessageStatus {
	defer group.recordMessage(responseMessage)

	if val, ok := group.Messages[string([]byte(responseMessage.requestHash[:]))]; ok {
		if val.MessageStatus == OnGoing {
			if group.IsRequestType(val.RequestName) {

				verifyFunc, ok2 := group.VerifyResponse[val.RequestName]

				if ok2 {
					val.MessageStatus, val.Responses = verifyFunc(responseMessage, val.Responses, val.AuthorizedUsers)
				} else {
					val.MessageStatus, val.Responses = group.DefaultVerifyResponse(responseMessage, val.Responses, val.AuthorizedUsers)
				}

				if val.MessageStatus == Success {
					delete(group.PendingRequests[val.RequestName], string([]byte(responseMessage.requestHash[:])))
					group.HandleRequestType(val.RequestName, val.RequestArgs...)
				} else if val.MessageStatus == Failure {
					delete(group.PendingRequests[val.RequestName], string([]byte(responseMessage.requestHash[:])))
				}

			} else {
				delete(group.PendingRequests[val.RequestName], string([]byte(responseMessage.requestHash[:])))
				val.MessageStatus = Failure
			}

			group.Messages[string([]byte(responseMessage.requestHash[:]))] = val
			return val.MessageStatus
		}
		return Failure
	} else {
		return Failure
	}
}

func (group *Group) GetPendingRequests(userID string) []RequestMessage {
	requestTypes := group.BasicGroup.GetRequestTypesForMember(userID)
	pending := make([]RequestMessage, 0)

	for _, val := range requestTypes {
		if requests, exists := group.PendingRequests[val]; exists {
			for key, _ := range requests {
				pending = append(pending, group.Messages[key].RequestMessage)
			}
		}
	}

	return pending
}
