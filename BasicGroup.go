package main

//All functions assert that the requested change actually is applied but the boolean return is to indicate if a change was required to finalize this request


type RequestHandler func(args ...interface{}) interface{}


type BasicGroup struct {
	Name        string //A name to be given to the group
	Description string //Descripes the main aims of the group

	Creator string //The creator of the group's username

	//Set that conatins all the userIDs of the members in a group (General assembly)
	Members map[string]struct{}

	//Stores the assignment of each member to some roles in the group, role members should be in the Members set
	//The first key is the role name
	//The second key is the userID
	Memberships map[string]map[string]struct{}

	//Stores which roles can handle each request type
	//First Key is the request type name
	//second key is the member userID
	Authorizations map[string]map[string]struct{}
}

//Checks if the user is a general member in the group
func (group *BasicGroup) IsMember(userID string) bool {
	_, exist := group.Members[userID];
	return exist
}

//Checks if the role is in the group
func (group *BasicGroup) IsRole(role string) bool {
	_, exist := group.Memberships[role];
	return exist
}

//checks if the user is a member in a certain role
func (group *BasicGroup) IsMemberInRole(userID string, role Role) bool {
	if group.IsRole(role){
		_, exist := group.Memberships[role][userID];
		return exist
	}
	return false
}

//Adds a user to the general Members list in the group
func (group *BasicGroup) AddMember(userID string) bool {
	if ! group.IsMember(userID){
		group.Members[userID] = struct{}
		return true
	}
	return false
}

//Adds a group member to an existing role
func (group *BasicGroup) AddMemberInRole(userID string, role string) bool {
	if !group.IsMemberInRole(userID, role) && group.IsRole(role){
		group.AddMember(userID)
		group.Memberships[role][userID] = struct{}
		return true
	}
	return false
}

//Adds a list of group members to an existing role
func (group *BasicGroup) AddInRole(userIDs []string, role string) bool {
	var change bool
	change  = false

	for userID := range UserIDs{
		change = change || group.AddMemberInRole(userID, role)
	}
	return true
}

//Remove a user from the group members
//and hence removes the user from any role
func (group *BasicGroup) RemoveMember(userID string) bool {
	if group.IsMember(userID){
		for key, _ := range group.Memberships {
			delete(group.Memberships[key], userID)
		}
		delete(group.Members, userID)
		return true;
	}

	return false
}


//Remove a user from a certain role in the group
//If the user is only in this role, the user will be removed from the general group members
func (group *BasicGroup) RemoveMemberInRole(userID string, role Role) bool {
	if group.IsMemberInRole(userID, role){
		delete(group.Memberships[role], userID)
		if len(group.GetMememberRoles(userID)) == 0{
			delete(group.Members, userID)
			}
		return true
	}

	return false
}

//adds a new empty role in the group
//If the role already exists nothing happens
func (group *BasicGroup) AddRole(role string)bool {
	if !group.IsRole(role){
		group.Memberships[role] = make(map[string]struct{}, 0)
		return true
	}
	return false
}

//removes an existing role in the group
//Therefore all members existing in only this role will be removed from the general members list
func (group *BasicGroup) RemoveRole(role string) bool {
	if group.IsRole(role){
		for key, val := range group.Memberships[role]{
			if len(group.GetMememberRoles(key, role)) == 1{
				delete(group.Members, key)
			}
		}
		delete(group.Memberships, role)
		return true
	}
	return false
}

//merges two different roles to create a new role in the group
//if the role already exists, it appends the members in the two roles in it
//if the keepOld flag is set false, old roles will removed
func (group *BasicGroup) MergeRoles(roleOne string, roleTwo string, newRole string, keepOld bool) {
	change := false

	change = change || group.AddRole(role) || group.AddInRole(group.GetMembersInRoles(roleOne, roleTwo), role)

	if !keepOld{
		change = change || group.RemoveRole(roleOne) || group.RemoveRole(roleTwo)
	}
	
	return change
}

//Add new request type in the group
//If one already exists, it oevrwrites it
//Only existing roles will be added to manage the new request type
func (group *BasicGroup) AddRequestType(requestType string, managingRoles []string, successFunction RequestHandler) bool {
}

//Remove an existing request type
func (group *BasicGroup) RemoveRequestType(requestType string) bool {
}

//Authorize a role to manage an existing type of requests
func (group *BasicGroup) AuthorizeRoleForRequestType(role string, requestType string) bool {
}

//Unauthorize a role from managing an existing type of requests
func (group *BasicGroup) UnauthorizeRoleFromRequestType(role string, requestType string) bool {
}


//Important getters

//Gets All the roles authorized to handle an existing request type
func (group *BasicGroup) GetRolesForRequestType(args ...interface{}) bool {
}


//Gets All the members belonging to at least one existing role in a set of roles
func (group *BasicGroup) GetMembersInRoles(args ...interface{}) bool {
}

//Gets All the members authorized to handle an existing request type
func (group *BasicGroup) GetMembersForRequestType(args ...interface{}) bool {
}


//Gets all the roles to which a member belongs
func (group BasicGroup) GetMememberRoles(userID string) []Role {
	var roles []Role

	for key, val := range group.Memberships {
		if val[userID] {
			roles = append(roles, key)
		}
	}

	return roles
}


//Gets all the request types that a member can handle
func (group BasicGroup) GetRequestTypesForMember(userRoles []Role) []MembershipRequestType {
	var authorizations []MembershipRequestType

	for key, roles := range group.Authorizations {
		Added := false
		for role, val := range roles {
			if val {
				for _, userRole := range userRoles {
					if userRole == role {
						authorizations = append(authorizations, key)
						Added = true
						break
					}
				}
				if Added {
					break
				}
			}
		}
	}

	return authorizations
}


//this is in the in other layer
//Get MemberPendingRequests
//Add Request
//Add a response
//GetRequestNumber
//GetResponseNumber
