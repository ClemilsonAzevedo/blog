package enums

type Role string

const (
	Anonymous  Role = "anonymous"
	Reader Role = "reader"
	Author = "author"
)
 
func (r *Role) VerifyRole() error {
	return nil
}
