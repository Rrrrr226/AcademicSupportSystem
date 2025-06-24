package oauth

import (
	"HelpStudent/internal/app/user/model/thirdPlat"
	"HelpStudent/pkg/utils/gen/xrandom"
	"fmt"
)

func GetRedirectUrl(feCallbackURL string, platform thirdPlat.Type, callbackURL string) (redirectURL string, mark string) {
	mark = xrandom.GetRandom(10, xrandom.RandAll)

	state := fmt.Sprintf("%s_%s", platform, mark)

	if PlatformExists(feCallbackURL, platform) {
		return PlatformEndpoint(feCallbackURL, platform).Redirect(callbackURL, state), mark
	}
	return "", ""
}
