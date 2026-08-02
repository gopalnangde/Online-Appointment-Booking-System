package constants

// Success messages
const (
	MsgRegistrationSuccess = "User registered successfully"
	MsgLoginSuccess        = "Login successful"
	MsgProfileFetched      = "Profile fetched successfully"
)

// Error messages
const (
	MsgEmailExists         = "Email already exists"
	MsgPhoneExists         = "Phone number already exists"
	MsgInvalidCredentials  = "Invalid email or password"
	MsgUserNotFound        = "User not found"
	MsgInvalidRequestBody  = "Invalid request body"
	MsgValidationFailed    = "Validation failed"
	MsgInternalServerError = "Something went wrong"
	MsgUnauthorized        = "Unauthorized access"
	MsgForbidden           = "You do not have permission to access this resource"
	MsgTokenMissing        = "Authorization token is required"
	MsgTokenInvalid        = "Invalid or expired token"
	MsgInvalidTokenFormat  = "Authorization header format must be Bearer {token}"
	MsgInvalidRole         = "Role must be either Customer or ServiceProvider"
)
