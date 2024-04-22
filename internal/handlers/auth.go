package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"responder/config"
	"responder/pkg/lc_api/auth"
	"time"
)

type authCodeCallbackHandler struct {
	config  *config.Config
	authAPi auth.LcAuthApi
}

func newAuthCodeCallbackHandler(config *config.Config, authApi auth.LcAuthApi) *authCodeCallbackHandler {
	return &authCodeCallbackHandler{
		config:  config,
		authAPi: authApi,
	}
}

func (nacch *authCodeCallbackHandler) Handle(w http.ResponseWriter, r *http.Request) {
	slog.Info("oAuth code received")
	code := r.URL.Query().Get("code")
	if code == "" {
		slog.Error("Cannot read code from path")
	}

	accessTokenReponse, err := nacch.authAPi.Exchange(auth.NewDefaultAccessTokenGrantRequest(code, nacch.config.ClientID, nacch.config.OauthConfig.ClientSecret, nacch.config.OauthConfig.RedirectURI))

	if err != nil {
		slog.Error("Cannot perform oauth flow; ", slog.Any("msg", err))
		w.Write([]byte("Retry the operation"))
		return
	}

	fmt.Println(accessTokenReponse.AccessToken)

	nacch.config.SetAccessToken(accessTokenReponse.AccessToken)
	nacch.config.SetRefreshToken(accessTokenReponse.RefreshToken)
	nacch.config.SetExpiresIn(accessTokenReponse.ExpiresIn)
	nacch.config.SetAccessTokenCreationTime(time.Now())


	w.Write([]byte("SUCCESS - you can close this window"))
}
