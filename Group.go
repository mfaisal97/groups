package groups

type Message struct {
	Request  interface{}
	Response []interface{}
}

type Group struct {
	BasicGroup

	Verify func(userIDs []string, signature interface{}, args ...interface{})

	Messages map[string]Message
}
