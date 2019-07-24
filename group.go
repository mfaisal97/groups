package main

type Group struct {
	//	Id          string
	Name        string
	Description string
	// Data        map[RequestType]Data

	Creator UserID

	Members        map[UserID]bool
	Memberships    map[Role]map[UserID]bool
	Authorizations map[MembershipRequestType]map[Role]bool
	//Membersauthorizations map[UserID][]RequestType
}

func (group *Group) IsMember(userID UserID) bool {
	return group.Members[userID]
}

func (group *Group) AddMember(userID UserID) bool {
	group.Members[userID] = true
	return true
}

func (group *Group) RemoveMember(userID UserID) bool {
	group.Members[userID] = true

	for _, val := range group.Memberships {
		if val[userID] {
			val[userID] = false
		}
	}

	return true
}

func (group *Group) IsMemberInRole(userID UserID, role Role) bool {
	return group.Memberships[role][userID]
}

func (group *Group) AddMemberInRole(userID UserID, role Role) bool {
	group.Members[userID] = true
	group.Memberships[role][userID] = true
	return true
}

func (group *Group) RemoveMemberInRole(userID UserID, role Role) bool {
	for key, _ := range group.Memberships[role] {
		if key == userID {
			group.Memberships[role][userID] = false
		}
	}

	if len(group.GetRoles(userID)) == 0 {
		group.RemoveMember(userID)
	}

	return true
}

func (group Group) GetRoles(userID UserID) []Role {
	var roles []Role

	for key, val := range group.Memberships {
		if val[userID] {
			roles = append(roles, key)
		}
	}

	return roles
}

func (group Group) GetAuthorizations(userRoles []Role) []MembershipRequestType {
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
