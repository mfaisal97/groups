package main

//stores all reqeusts and replies and verify sender.
type DataBase struct {
	Groups Group

	GroupCreatationMessages GroupCreatationMessage
	MembershipsMessages     MembershipMessage

	Users User
}
