package data

type Role struct {
	ID   int32
	Name string
}

var (
	MemberRole = &Role{
		ID:   0,
		Name: "Member",
	}
	OwnerRole = &Role{
		ID:   1,
		Name: "Owner",
	}
	AdministratorRole = &Role{
		ID:   2,
		Name: "Admin",
	}
	Roles = []*Role{
		MemberRole,
		OwnerRole,
		AdministratorRole,
	}
)

const (
	MembershipAccepted = "Accepted"
	MembershipPending  = "Pending"
	MembershipRejected = "Rejected"
)
