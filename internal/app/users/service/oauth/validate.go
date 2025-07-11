package oauth

import (
	"HelpStudent/internal/app/users/model/thirdPlat"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

// Validate 返回第三方平台用户unique id
func Validate(feCallbackURL string, platform thirdPlat.Type, code string, state string) (id string, attr datatypes.JSON, err error) {
	if PlatformExists(feCallbackURL, platform) {
		return PlatformEndpoint(feCallbackURL, platform).Validate(code, state)
	}
	return "", nil, errors.New("platform not supported")
}
