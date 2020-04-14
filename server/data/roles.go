package data

type Role struct {
	ID              int32
	Name            string
	readOnlyProject bool
}

var (
	MemberRole = &Role{
		ID:              0,
		Name:            "Member",
		readOnlyProject: true,
	}
	OwnerRole = &Role{
		ID:              1,
		Name:            "Owner",
		readOnlyProject: false,
	}
	AdministratorRole = &Role{
		ID:              2,
		Name:            "Admin",
		readOnlyProject: false,
	}
	PublicRole = &Role{
		// NOTE This is never persisted.
		Name:            "Public",
		readOnlyProject: true,
	}
	Roles = []*Role{
		MemberRole,
		OwnerRole,
		AdministratorRole,
	}
)

func (r *Role) IsProjectReadOnly() bool {
	return r.readOnlyProject
}

const (
	MembershipAccepted = "Accepted"
	MembershipPending  = "Pending"
	MembershipRejected = "Rejected"
)
