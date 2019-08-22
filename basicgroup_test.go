package main

import (
	"fmt"
	"testing"
)

func sum(args ...interface{}) interface{} {
	return args[0].(int) + args[1].(int)
}

func TestBasicgroup(t *testing.T) {

	passedTestsNum := 0
	totalTests := 10
	printObjects := true

	//Initializing import
	group := CreateNewBasicGroup("groupName", "description", "creator")
	group.SetCreatorName("AnotherName")
	if group.GetName() == "groupName" && group.GetDescription() == "description" && group.GetCreatorName() == "creator" {
		passedTestsNum++
		fmt.Println("Testing basic group creation passed:\t name - description - creatorname \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//Adding members in non existing roles
	toAdmins := make([]string, 0, 3)
	toAdmins = append(toAdmins, "John")
	toAdmins = append(toAdmins, "Mike")
	toAdmins = append(toAdmins, "Mike")
	toNonmembers := make([]string, 0)
	group.AddMemberInRole("mickey", "Admins")
	group.AddInRole(toAdmins, "Admins")
	group.AddInRole(toNonmembers, "NonAdmins")
	managingRoles := make([]string, 0, 1)
	managingRoles = append(managingRoles, "NonAdmins")
	group.AddRequestType("requestType", managingRoles, sum)
	if !group.IsMember("mickey") && !group.IsRole("Admins") && len(group.GetRolesForRequestType("requestType")) == 0 && len(group.GetRolesForRequestType("NorequestType")) == 0 {
		passedTestsNum++
		fmt.Println("Adding members in non existing roles passed:\t mickey/John,Mike,Mike >> Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//removing members from non existing roles
	group.RemoveMemberInRole("mickey", "Amdins")
	if !group.IsMember("mickey") && !group.IsRole("Admins") {
		passedTestsNum++
		fmt.Println("removing members from non existing roles passed:\t mickey/John,Mike,Mike >> Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//creating new roles and adding members
	group.AddRole("Admins")
	group.AddMemberInRole("mickey", "Admins")
	group.AddInRole(toAdmins, "Admins")
	group.AddInRole(toNonmembers, "NonAdmins")
	if group.IsMember("mickey") && group.IsMember("John") && group.IsMember("Mike") {
		passedTestsNum++
		fmt.Println("creating new roles and adding members passed:\t mickey/John,Mike,Mike >> Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//recreating existing role and authorizing role to manage a request
	group.AddRole("Admins")
	group.AuthorizeRoleForRequestType("Admins", "requestType")
	if group.IsRole("Admins") && group.IsAuthorizedRole("Admins", "requestType") && len(group.GetMembersForRequestType("requestType")) == 3 {
		passedTestsNum++
		fmt.Println("recreating existing role passed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//trying removing members and request types for a member
	group.RemoveMemberInRole("Mike", "Admins")
	if group.IsRole("Admins") && !group.IsMember("Mike") && !group.IsMemberInRole("Mike", "Admins") && len(group.GetRequestTypesForMember("Mike")) == 0 && group.IsAuthorizedMember("mickey", "requestType") {
		passedTestsNum++
		fmt.Println("trying removing members passed:\t Mike From Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//checking geting MembersInRoles
	group.AddRole("NonMembers")
	group.AddInRole(toNonmembers, "NonMembers")
	group.AddMemberInRole("noOne", "NonMembers")
	if group.IsRole("NonMembers") && len(group.GetMembersInRoles("NonMembers")) == 1 && len(group.GetMembersInRoles("NonMembers", "Admins")) == 3 {
		passedTestsNum++
		fmt.Println("checking geting MembersInRoles passed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//trying removing members that belong to one role
	group.RemoveMemberInRole("noOne", "NonMembers")
	group.AddRole("TwoPersonRole")
	group.AddMemberInRole("FirstPerson", "TwoPersonRole")
	group.AddMemberInRole("SecondPerson", "TwoPersonRole")
	group.AddMemberInRole("FirstPerson", "Admins")
	if group.IsRole("NonMembers") && len(group.GetMembersInRoles("NonMembers")) == 0 && len(group.GetMembersInRoles("NonMembers", "Admins")) == 3 {
		passedTestsNum++
		fmt.Println("trying removing members that belong to one role passed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//removing role that has a member belonging only to it
	group.RemoveRole("TwoPersonRole")
	if !group.IsRole("TwoPersonRole") && group.IsMemberInRole("FirstPerson", "Admins") && !group.IsMemberInRole("FirstPerson", "TwoPersonRole") && !group.IsMemberInRole("SecondPerson", "TwoPersonRole") {
		passedTestsNum++
		fmt.Println("removing role that has a member belonging only to itpassed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//merging roles and checking authorization for non existing request type
	group.AddMemberInRole("userID", "NonMembers")
	group.MergeRoles("TwoPersonRole", "NonMembers", "newRole", true)
	group.MergeRoles("TwoPersonRole", "NonMembers", "newRoleAfterDeletion", false)
	group.AuthorizeRoleForRequestType("newRole", "requestType")
	if !group.IsRole("TwoPersonRole") && !group.IsRole("NonMembers") && !group.IsMemberInRole("userID", "NonMembers") && group.IsMemberInRole("userID", "newRole") && group.IsMemberInRole("userID", "newRoleAfterDeletion") && !group.IsAuthorizedMember("notUser", "requestType") && !group.IsAuthorizedRole("NotRole", "NotrequestType") && len(group.GetRequestTypesForMember("John")) == 1 {
		passedTestsNum++
		fmt.Println("removing role that has a member belonging only to itpassed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	group.removeMember("mickey")
	group.UnauthorizeRoleFromRequestType("newRole", "requestType")
	if !group.IsMemberInRole("mickey", "Admins") && !group.IsMember("mickey") && len(group.GetGroupMembers()) == 3 {
		passedTestsNum++
		fmt.Println("removing role that has a member belonging only to itpassed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//deleting all roles
	group.RemoveRole("Admins")
	group.RemoveRole("newRole")
	group.RemoveRole("newRoleAfterDeletion")
	if len(group.GetGroupMembers()) == 0 && !group.IsRole("Admins") && !group.IsRole("newRole") && !group.IsRole("newRoleAfterDeletion") {
		passedTestsNum++
		fmt.Println("removing role that has a member belonging only to itpassed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	// creating many roles with one containing all the group members
	group.removeMember("Noone")
	group.AddRole("Admins")
	group.AddRole("Members")
	group.AddRole("SuperAdmins")
	group.AddInRole(toAdmins, "Admins")
	group.AddMemberInRole("mickey", "SuperAdmins")
	group.AddMemberInRole("Mike", "Members")
	res, _ := group.HandleRequestType("requestType", 4, 5)
	group.UnauthorizeRoleFromRequestType("newRole", "requestType")
	if len(group.GetGroupMembers()) == 3 && group.IsRole("Admins") && group.IsRole("Members") && group.IsRole("SuperAdmins") && res == 9 {
		passedTestsNum++
		fmt.Println("removing role that has a member belonging only to itpassed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

	//removing the role that contains all the group members
	group.RemoveRole("Admins")
	group.RemoveRequestType("requestType")
	_, done := group.HandleRequestType("requestType", 4, 5)
	if len(group.GetGroupMembers()) == 2 && !group.IsRole("Admins") && group.IsRole("Members") && group.IsRole("SuperAdmins") && !done {
		passedTestsNum++
		fmt.Println("removing role that has a member belonging only to itpassed:\t role:Admins \t\t", passedTestsNum, "of", totalTests)
	}
	if printObjects {
		fmt.Println(group)
	}

}
