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
	PublicRole = &Role{
		// NOTE This is never persisted.
		Name:            "Public",
		readOnlyProject: true,
	}
	Roles = []*Role{
		MemberRole,
		OwnerRole,
		PublicRole,
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
