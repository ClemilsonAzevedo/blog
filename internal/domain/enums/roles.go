package enums

type Role int

const (
	RoleAnonymous Role = iota
	RoleReader
	RoleAuthor
)

func (r *Role) VerifyRole() error {
	return nil
}
