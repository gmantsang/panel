package radios

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

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
	radios[0].LastTested = strings.Split(radios[0].LastTested, "T")[0]

	ctx := pongo2.Context{}
	ctx["radio"] = radios[0]
	ctx["action"] = "edit"

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	utils.AddAuthContext(session, ctx, handler.Config)
	utils.AddMetaContext(ctx, handler.Meta)

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
		return
	}
	radio := parseForm(r)
	ok, reason := isValid(radio)
	if !ok {
		http.Error(w, reason, http.StatusBadRequest)
		return
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

func parseForm(r *http.Request) api.Radio {
	return api.Radio{
		Name:       r.Form.Get("Name"),
		URL:        r.Form.Get("URL"),
		Category:   r.Form.Get("Category"),
		Genre:      r.Form.Get("Genre"),
		Country:    r.Form.Get("Country"),
		LastTested: r.Form.Get("LastTested"),
		State:      r.Form.Get("State"),
	}
}

func isValid(r api.Radio) (bool, string) {
	if r.Name == "" {
		return false, "name"
	}
	if r.URL == "" {
		return false, "url"
	}
	if r.Category == "" {
		return false, "category"
	}
	if r.Genre == "" {
		return false, "genre"
	}
	if r.Country == "" {
		return false, "country"
	}
	if r.LastTested == "" {
		return false, "last_tested"
	}
	if r.State == "" {
		return false, "state"
	}
	return true, ""
}
