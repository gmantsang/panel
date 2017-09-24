package radios

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/flosch/pongo2"

	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/dabbotorg/radio-api/api"
)

// ViewBreak renders a confirmation page to move a radio into the broken state
func (handler *Handler) ViewBreak(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)

	ctx := pongo2.Context{}
	ctx["name"] = name
	ctx["action"] = "break"

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	utils.AddAuthContext(session, ctx, handler.Config)

	err = handler.Templates.Confirm.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, utils.TemplateFailed(err), http.StatusInternalServerError)
		return
	}
}

// Break handles a confirmation to flag a radio as broken
func (handler *Handler) Break(w http.ResponseWriter, r *http.Request) {
	handler.changeState(w, r, "BROKEN")
	http.Redirect(w, r, "/radios/list/valid", http.StatusTemporaryRedirect)
}

func (handler *Handler) changeState(w http.ResponseWriter, r *http.Request, state string) {
	name := path.Base(r.URL.Path)

	url := utils.GetRadioURL(handler.Config.APIUrl, name)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, utils.APIError(err, resp), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, utils.IOError(err), http.StatusInternalServerError)
		return
	}

	var radios []api.Radio
	err = json.Unmarshal(body, &radios)
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}
	radios[0].State = state
	raw, err := json.Marshal(radios[0])
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}
	buf := bytes.NewBuffer(raw)

	req, _ := http.NewRequest(http.MethodPatch, url, buf)
	req.Header.Set("Authorization", handler.Config.APIToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, utils.APIError(err, resp), http.StatusInternalServerError)
		return
	}
}
