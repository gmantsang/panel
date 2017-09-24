package radios

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/dabbotorg/radio-api/api"
	"github.com/flosch/pongo2"
)

// ViewEdit renders the radio-edit template
func (handler *Handler) ViewEdit(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)

	url := utils.GetRadioURL(handler.Config.APIUrl, name)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, utils.APIError(err, resp), http.StatusInternalServerError)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, utils.IOError(err), http.StatusInternalServerError)
		return
	}

	var radios []api.Radio
	err = json.Unmarshal(buf, &radios)
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}

	ctx := pongo2.Context{}
	ctx["radio"] = radios[0]

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	utils.AddAuthContext(session, ctx, handler.Config)

	err = handler.Templates.Edit.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, utils.TemplateFailed(err), http.StatusInternalServerError)
	}
}

// Edit handles a POST request to edit a radio
func (handler *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	radio := api.Radio{
		Name:       r.Form.Get("Name"),
		URL:        r.Form.Get("URL"),
		Category:   r.Form.Get("Category"),
		Genre:      r.Form.Get("Genre"),
		Country:    r.Form.Get("Country"),
		LastTested: r.Form.Get("LastTested"),
		State:      r.Form.Get("State"),
	}

	buf, err := json.Marshal(radio)
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}

	url := utils.GetRadioURL(handler.Config.APIUrl, name)
	reader := bytes.NewReader(buf)
	req, _ := http.NewRequest("PATCH", url, reader)
	req.Header.Set("Authorization", handler.Config.APIToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, utils.APIError(err, resp), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/radios/list/valid", http.StatusTemporaryRedirect)
}
