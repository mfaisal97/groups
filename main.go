package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

//global variables
var data DataBase
var userID string
var memberID string
var memberRole string
var groupName string
var members string
var admins string
var creator string
var stringAuth string
var accepted bool
var requestNum int

func ResetDataBase(c *cli.Context) {
	data.ResetDataBase()
}

func CreateUser(c *cli.Context) {
	if c.NArg() != 1 {
		_ = errors.New("please give the following argument: UserID")
		return
	}
	userid := c.Args().First()
	_ = CreateNewUser(UserID(userid))
}

func CreateGroup(c *cli.Context) {
	group := Group{}
	group.Creator = UserID(creator)
	group.Name = groupName
	group.Description = "Just Testing Groups"

	group.Memberships = make(map[Role]map[UserID]bool)
	group.Memberships["member"] = make(map[UserID]bool)
	group.Members = make(map[UserID]bool)
	for _, val := range strings.Split(members, "-") {
		group.Members[UserID(val)] = true
		group.Memberships["member"][UserID(val)] = true
	}

	group.Memberships["admin"] = make(map[UserID]bool)
	for _, val := range strings.Split(admins, "-") {
		group.Members[UserID(val)] = true
		group.Memberships["admin"][UserID(val)] = true
	}

	group.Authorizations = make(map[MembershipRequestType]map[Role]bool)
	group.Authorizations[Join] = make(map[Role]bool)
	group.Authorizations[Join]["admin"] = true
	group.Authorizations[Remove] = make(map[Role]bool)
	group.Authorizations[Remove]["Admin"] = true

	user := User{}
	user.LoadData(UserID(userID))
	user.CreateGroup(group)

}

func JoinGroup(c *cli.Context) {
	user := User{}
	user.LoadData(UserID(userID))
	user.SendMembershipRequest(Join, groupName, UserID(memberID), Role(memberRole))
}

func RemoveMember(c *cli.Context) {
	user := User{}
	user.LoadData(UserID(userID))
	user.SendMembershipRequest(Remove, groupName, UserID(memberID), Role(memberRole))
}

func ReplyRequest(c *cli.Context) {
	user := User{}
	user.LoadData(UserID(userID))
	user.SendMembershipReesponse(int32(requestNum), accepted)
}

func GetRequests(c *cli.Context) {
	user := User{}
	user.LoadData(UserID(userID))
	user.GetPendingRequests(groupName)
}

func main() {
	data = InitializeDataBase()
	data.LoadData("Data/data.json")
	fmt.Println("Starting\t", data)

	cliApp := cli.NewApp()
	cliApp.Name = "Groups-Management"
	cliApp.Version = "0.0"
	cliApp.Commands = []cli.Command{
		{
			Name:   "reset",
			Usage:  "Removed the previous data.json file",
			Action: ResetDataBase,
		},
		{
			Name:   "createUser",
			Usage:  "Creates a new user in the data base",
			Action: CreateUser,
		},
		{
			Name:   "createGroup",
			Usage:  "Creates a new Group in the database",
			Action: CreateGroup,
		},
		{
			Name:   "joinGroup",
			Usage:  "Sends a command to add a user in the group to a specific role",
			Action: JoinGroup,
		},
		{
			Name:   "removeMember",
			Usage:  "Sends a command to remove a user in the group from a specific role ",
			Action: RemoveMember,
		},
		{
			Name:   "replyRequest",
			Usage:  "Sends a command to reply a specific request in the database",
			Action: ReplyRequest,
		},
		{
			Name:   "GetRequests",
			Usage:  "Get all possible requests to reply to in a group",
			Action: GetRequests,
		},
	}

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "user, u",
			Usage:       "userID the member initiating the command",
			Destination: &userID,
		},
		cli.StringFlag{
			Name:        "group, g",
			Usage:       "name of the mentioned group",
			Destination: &groupName,
		},
		cli.StringFlag{
			Name:        "members, m",
			Usage:       "UserIDs of the group members",
			Destination: &members,
		},
		cli.StringFlag{
			Name:        "admins, a",
			Usage:       "UserIDs of the group admins",
			Destination: &admins,
		},
		cli.StringFlag{
			Name:        "creator, c",
			Usage:       "UserID of the group creator",
			Destination: &creator,
		},
		// cli.StringFlag{
		// 	Name:        "authorize, z",
		// 	Usage:       "Pairs of MembershipRequestTypes and roles",
		// 	Destination: &stringAuth,
		// },
		cli.StringFlag{
			Name:        "member",
			Usage:       "UserID of the group member the request is affecting",
			Destination: &memberID,
		},
		cli.StringFlag{
			Name:        "role",
			Usage:       "the role of the member to be changed",
			Destination: &memberRole,
		},
		cli.BoolFlag{
			Name:        "reply, r",
			Usage:       "UserID of the group creator",
			Destination: &accepted,
		},
		cli.IntFlag{
			Name:        "requestNum, n",
			Usage:       "The request number of the request to reply to",
			Destination: &requestNum,
		},
	}

	cliApp.After = func(c *cli.Context) error {
		return nil
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	data.SaveData("Data/data.json")
	fmt.Println("Ending\t", data, "\n\n")
}
