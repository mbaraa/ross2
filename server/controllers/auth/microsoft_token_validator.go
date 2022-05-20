package auth

// Microsoft JWT validation stuff
// 😉e10

type MicrosoftJWTValidator struct{}

func NewMicrosoftJWTValidator() *MicrosoftJWTValidator {
	return new(MicrosoftJWTValidator)
}

func (m *MicrosoftJWTValidator) Validate(token string) error {
	return nil
}
