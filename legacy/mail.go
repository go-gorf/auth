package auth

type AuthEmailVerificationBackend interface {
	PrepareMail(user *User) error
	SendMail() error
}
