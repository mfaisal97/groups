package groups

type UserInfo struct {
	UserID     string
	PublicKey  interface{}
	PrivateKey interface{}
}

type GroupMemberInfo struct {
	PendingRequests []Message
	SignRequest     func(request Request, userInfo UserInfo) interface{}
	SignResponse    func(response Response, userInfo UserInfo, userIDs []string) interface{}
}

type Request struct {
	RequestName      string
	RandomIdentifier int
	RequestArgs      []interface{}
}

type RequestMessage struct {
	Request
	Signature interface{}
}

type Response struct {
	requestHash [32]byte
	Answer      bool
}

type ResponseMessage struct {
	Response
	Signature interface{}
}

type Message struct {
	RequestMessage
	Responses       map[string]ResponseMessage
	AuthorizedUsers []string
	MessageStatus
}

type MessageStatus int

const (
	OnGoing MessageStatus = 0
	Success MessageStatus = 1
	Failure MessageStatus = 2
)

type RequestStatus int

const (
	RequestConfirmed RequestStatus = 0
	RequestReceived  RequestStatus = 1
	RequestFailed    RequestStatus = 2
)
