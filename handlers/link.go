package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
)

type (
	// LinkHandler handles /link/discord
	LinkHandler struct {
		Store        sessions.Store
		ClientID     string
		ClientSecret string
	}

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

func (handler *LinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		handler.HandleGet(w, r)
		return
	}
	handler.HandleCallback(w, r, code)
}

// HandleGet handles the initial GET /link/discord, which will redirect to Discord's OAuth
func (handler *LinkHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	session, err := handler.Store.Get(r, SessionName)
	if err != nil {
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
		return
	}
	if session.Values["authed"] == "true" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, DiscordAuthURL(handler.ClientID), http.StatusTemporaryRedirect)
}

// HandleCallback handles the OAuth callback from Discord
func (handler *LinkHandler) HandleCallback(w http.ResponseWriter, r *http.Request, code string) {
	errorQuery := r.URL.Query().Get("error")
	if errorQuery != "" {
		http.Error(w, DiscordError(errorQuery), http.StatusBadRequest)
		return
	}

	url := DiscordTokenURL(handler.ClientID, handler.ClientSecret, code)
	tokenResponse, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		http.Error(w, TokenExchangeFailed(err), http.StatusInternalServerError)
	}

	buf, err := ioutil.ReadAll(tokenResponse.Body)
	if err != nil {
		http.Error(w, IOError(err), http.StatusInternalServerError)
		return
	}
	var tokenData bearerResponse
	err = json.Unmarshal(buf, &tokenData)
	if err != nil {
		http.Error(w, JSONError(err), http.StatusInternalServerError)
		return
	}

	req, _ := http.NewRequest(http.MethodGet, DiscordMeURL, nil)
	token := fmt.Sprintf("Bearer %s", tokenData.Token)
	req.Header.Set("Authorization", token)
	meResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, DiscordError(err.Error()), http.StatusInternalServerError)
	}

	buf, err = ioutil.ReadAll(meResponse.Body)
	if err != nil {
		http.Error(w, IOError(err), http.StatusInternalServerError)
		return
	}
	var userData userResponse
	err = json.Unmarshal(buf, &userData)
	if err != nil {
		http.Error(w, JSONError(err), http.StatusInternalServerError)
		return
	}

	session, err := handler.Store.Get(r, SessionName)
	if err != nil {
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
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
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
