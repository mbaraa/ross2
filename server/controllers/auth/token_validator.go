package auth

type JWTTokenValidator interface {
	Validate(token string) error
}
