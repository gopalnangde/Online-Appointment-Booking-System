package constants

// User role constants
const (
	RoleCustomer        = "Customer"
	RoleServiceProvider = "ServiceProvider"
)

// ValidRoles contains all allowed role values
var ValidRoles = []string{RoleCustomer, RoleServiceProvider}

// IsValidRole checks if the given role is one of the allowed values
func IsValidRole(role string) bool {
	for _, r := range ValidRoles {
		if r == role {
			return true
		}
	}
	return false
}
