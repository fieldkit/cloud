package data

type Role struct {
	ID              int32
	Name            string
	persisted       bool
	readOnlyProject bool
}

var (
	MemberRole = &Role{
		ID:              0,
		Name:            "Member",
		readOnlyProject: true,
	}
	AdministratorRole = &Role{
		ID:              1,
		Name:            "Administrator",
		readOnlyProject: false,
	}
	PublicRole = &Role{
		// NOTE This is never persisted.
		Name:            "Public",
		readOnlyProject: true,
	}
	AvailableRoles = []*Role{
		MemberRole,
		AdministratorRole,
	}
	Roles = []*Role{
		MemberRole,
		AdministratorRole,
		PublicRole,
	}
)

func LookupRole(id int32) *Role {
	for _, r := range Roles {
		if r.ID == id {
			return r
		}
	}
	return nil
}

func (r *Role) IsProjectReadOnly() bool {
	return r.readOnlyProject
}

const (
	MembershipAccepted = "Accepted"
	MembershipPending  = "Pending"
	MembershipRejected = "Rejected"
)
