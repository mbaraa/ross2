package multiavatar

import (
	"fmt"

	"github.com/google/uuid"
)

func GetAvatarURL() string {
	return fmt.Sprintf("https://api.multiavatar.com/%s.svg", uuid.NewString())
}
