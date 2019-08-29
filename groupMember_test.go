package groups

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/group/edwards25519"
	"go.dedis.ch/kyber/sign/anon"
)

type testSignature struct {
	userID string
	sign   []byte
}

type testUsers struct {
	PublicKeys  map[string]kyber.Point
	PrivateKeys map[string]kyber.Scalar
}

func (testUsers *testUsers) verifyRequestTest(request Request, signature interface{}) RequestStatus {
	if publickKey, exist := testUsers.PublicKeys[signature.(testSignature).userID]; exist {
		suite := edwards25519.NewBlakeSHA256Ed25519()
		pubSet := make([]kyber.Point, 1)
		pubSet[0] = publickKey

		hash := GenerateHash(request)

		tag, err := anon.Verify(suite, []byte(hash[:]), anon.Set(pubSet), nil, signature.(testSignature).sign)
		if err != nil {
			panic(err.Error())
		}
		if tag == nil || len(tag) != 0 {
			panic("Verify returned wrong tag")
		}
		return RequestConfirmed
	}
	return RequestFailed

}

func (testUsers *testUsers) verifyResponseTest(responseMessage ResponseMessage, previousResponses map[string]ResponseMessage, userIDs []string) (MessageStatus, map[string]ResponseMessage) {
	responsehash := GenerateHash(responseMessage.Response)

	suite := edwards25519.NewBlakeSHA256Ed25519()
	var pubSet []kyber.Point
	for _, user := range userIDs {
		if publickKey, exists := testUsers.PublicKeys[user]; exists {
			pubSet = append(pubSet, publickKey)
		}
	}

	tag, err := anon.Verify(suite, []byte(responsehash[:]), anon.Set(pubSet), []byte(responseMessage.requestHash[:]), responseMessage.Signature.(testSignature).sign)
	if err != nil {
		panic(err.Error())
	}
	if tag == nil || len(tag) != suite.PointLen() {
		panic("Verify returned invalid tag")
	}

	if _, signExists := previousResponses[string(tag)]; signExists {
		return OnGoing, previousResponses
	}
	previousResponses[string(tag)] = responseMessage

	approved := 0
	for _, val := range previousResponses {
		if val.Answer {
			approved++
		}
	}

	minimumToWin := int(math.Ceil(float64(len(pubSet)) / 2.0))

	if approved >= minimumToWin {
		return Success, previousResponses
	} else if (len(previousResponses) - approved) >= minimumToWin {
		return Failure, previousResponses
	} else {
		return OnGoing, previousResponses
	}
}

func (testUsers *testUsers) signRequestTest(request Request, userInfo UserInfo) interface{} {
	suite := edwards25519.NewBlakeSHA256Ed25519()

	pubSet := make([]kyber.Point, 1)
	pubSet[0] = userInfo.PublicKey.(kyber.Point)

	sign := testSignature{}
	//sign.userID = userInfo.UserID
	hash := GenerateHash(request)
	sign.sign = anon.Sign(suite, []byte(hash[:]), anon.Set(pubSet), nil, 0, userInfo.PrivateKey.(kyber.Scalar))

	return sign
}

func (testUsers *testUsers) signResponseTest(response Response, userInfo UserInfo, userIDs []string) interface{} {
	var pubSet []kyber.Point
	for _, user := range userIDs {
		if publickKey, exists := testUsers.PublicKeys[user]; exists {
			pubSet = append(pubSet, publickKey)
		}
	}

	suite := edwards25519.NewBlakeSHA256Ed25519()
	responseHash := GenerateHash(response)

	index := 0
	for ind, val := range pubSet {
		if reflect.DeepEqual(userInfo.PublicKey, val) {
			index = ind
			break
		}
	}

	sign := testSignature{}
	sign.sign = anon.Sign(suite, []byte(responseHash[:]), anon.Set(pubSet), []byte(response.requestHash[:]), index, userInfo.PrivateKey.(kyber.Scalar))
	return sign
}

func (testUsers *testUsers) addUser(userID string, publicKey kyber.Point, privateKey kyber.Scalar) {
	testUsers.PublicKeys[userID] = publicKey
	testUsers.PrivateKeys[userID] = privateKey
}

func (testUsers *testUsers) removeUser(userID string) {
	delete(testUsers.PublicKeys, userID)
	delete(testUsers.PrivateKeys, userID)
}

func (testUsers *testUsers) createTestUser(userID string) GroupMember {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	privateKey := suite.Scalar().Pick(suite.RandomStream())
	publicKey := suite.Point().Mul(privateKey, nil)
	member := CreateMember(userID, privateKey, publicKey)

	testUsers.addUser(userID, publicKey, privateKey)

	return member
}

func TestGroupMembergroup(t *testing.T) {
	fmt.Println("Starting Test Two: \t\t Testing GroupMember and Group Classes")

	// passedTestsNum := 0
	// currentTestCase := 0
	// totalTests := 10
	// printObjects := true

}
