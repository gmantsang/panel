package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dabbotorg/panel/handlers/utils"
)

type (
	bearerResponse struct {
		Token   string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}
	userResponse struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar_id"`
	}
)

// StartFlow handles the initial GET /link/discord, which will redirect to Discord's OAuth
func (handler *Handler) StartFlow(w http.ResponseWriter, r *http.Request) {
	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	if session.Values["authed"] == "true" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, utils.DiscordAuthURL(handler.ClientID), http.StatusTemporaryRedirect)
}

// EndFlow handles the OAuth callback from Discord
func (handler *Handler) EndFlow(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	errorQuery := r.URL.Query().Get("error")
	if errorQuery != "" {
		http.Error(w, utils.DiscordError(errorQuery), http.StatusBadRequest)
		return
	}

	url := utils.DiscordTokenURL(handler.ClientID, handler.ClientSecret, code)
	tokenResponse, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		http.Error(w, utils.TokenExchangeFailed(err), http.StatusInternalServerError)
	}

	buf, err := ioutil.ReadAll(tokenResponse.Body)
	if err != nil {
		http.Error(w, utils.IOError(err), http.StatusInternalServerError)
		return
	}
	var tokenData bearerResponse
	err = json.Unmarshal(buf, &tokenData)
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}

	req, _ := http.NewRequest(http.MethodGet, utils.DiscordMeURL, nil)
	token := fmt.Sprintf("Bearer %s", tokenData.Token)
	req.Header.Set("Authorization", token)
	meResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, utils.DiscordError(err.Error()), http.StatusInternalServerError)
	}

	buf, err = ioutil.ReadAll(meResponse.Body)
	if err != nil {
		http.Error(w, utils.IOError(err), http.StatusInternalServerError)
		return
	}
	var userData userResponse
	err = json.Unmarshal(buf, &userData)
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}

	session.Values["access_token"] = tokenData.Token
	session.Values["refresh_token"] = tokenData.Refresh
	session.Values["authed"] = "true"
	session.Values["id"] = userData.ID
	session.Values["username"] = userData.Username
	session.Values["discrim"] = userData.Discriminator

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
