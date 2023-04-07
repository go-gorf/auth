package auth

type AuthState bool

var activateUserOnCreation AuthState = true
var skipEmailVerification AuthState = false
var defaultAdminStatus AuthState = false

// modify the setting like this in project settings
// auth.AuthSettings.NewUserState = auth.AuthState(true)
var AuthSettings = struct {
	NewUserState             AuthState
	NewUserAdminState        AuthState
	EmailVerification        AuthState
	EmailVerificationBackend AuthEmailVerificationBackend
}{
	NewUserState:      activateUserOnCreation,
	EmailVerification: skipEmailVerification,
	NewUserAdminState: defaultAdminStatus,
}
