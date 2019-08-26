package groups

type Group struct {
	BasicGroup

	Messages    map[string]Message
	AllMessages []interface{}
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
					return RequestReceived
				}
			}
		}
	}
	return RequestFailed
}

func (group *Group) AddResponse(responseMessage ResponseMessage) MessageStatus {
	defer group.recordMessage(responseMessage)

	if val, ok := group.Messages[string([]byte(responseMessage.requestHash[:]))]; ok && val.MessageStatus == OnGoing {
		if group.IsRequestType(val.RequestName) {

			verifyFunc, ok2 := group.VerifyResponse[val.RequestName]

			if ok2 {
				val.MessageStatus, val.Responses = verifyFunc(responseMessage, val.Responses, val.AuthorizedUsers)
			} else {
				val.MessageStatus, val.Responses = group.DefaultVerifyResponse(responseMessage, val.Responses, val.AuthorizedUsers)
			}

			if val.MessageStatus == Success {
				group.HandleRequestType(val.RequestName, val.RequestArgs...)
			}

		} else {
			val.MessageStatus = Failure
		}

		group.Messages[string([]byte(responseMessage.requestHash[:]))] = val
		return val.MessageStatus
	} else {
		return Failure
	}
}
