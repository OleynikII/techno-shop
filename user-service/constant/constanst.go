package constant

const (
	// User
	AdminRoleName      string = "admin"
	DefaultRoleName    string = "default"
	DefaultUserName    string = "admin"
	RedisOtpDefaultKey string = "otp"

	// Claims
	AuthorizationHeaderKey string = "Authorization"
	UserIdKey              string = "UserId"
	RolesKey               string = "Roles"
	ExpireTimeKey          string = "Exp"

	// JWT
	AccessTokenCookieName  string = "access_token"
	RefreshTokenCookieName string = "refresh_token"
)
