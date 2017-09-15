package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"path"
	"strconv"

	"github.com/dabbotorg/panel/config"
	"github.com/dabbotorg/radio-api/api"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

// RadioHandler manages radio stations
type RadioHandler struct {
	Store        sessions.Store
	ListTemplate *pongo2.Template
	EditTemplate *pongo2.Template
	Config       config.Config
}

func (handler *RadioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := path.Base(r.URL.Path)
	if action == "list" {
		handler.ListRadios(w, r, "VALID")
	} else if action == "escrow" {
		handler.ListRadios(w, r, "ESCROW")
	} else if action == "broken" {
		handler.ListRadios(w, r, "BROKEN")
	} else if action == "edit" {
		handler.EditRadio(w, r)
	} else {
		http.Error(w, "", http.StatusNotFound)
	}
}

// ListRadios lists radios
func (handler *RadioHandler) ListRadios(w http.ResponseWriter, r *http.Request, state string) {
	query := r.URL.Query()

	limit := query.Get("limit")
	if limit == "" {
		limit = fmt.Sprintf("100")
	}
	offset := query.Get("offset")
	if offset == "" {
		offset = fmt.Sprintf("0")
	}
	offsetInt, err := strconv.ParseInt(offset, 10, 32)
	if err != nil || offsetInt < 0 {
		http.Error(w, "offset must be an integer >= 0", http.StatusBadRequest)
		return
	}

	url := ListRadiosURL(handler.Config.APIUrl, limit, offset, state)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, APIError(err), http.StatusInternalServerError)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, IOError(err), http.StatusInternalServerError)
		return
	}

	var radios []api.Radio
	err = json.Unmarshal(buf, &radios)
	if err != nil {
		http.Error(w, JSONError(err), http.StatusInternalServerError)
		return
	}

	ctx := pongo2.Context{}
	ctx["radios"] = radios
	ctx["count"] = len(radios)
	ctx["page"] = math.Floor(float64((offsetInt + 100)) / float64(100))
	ctx["state"] = state

	session, err := handler.Store.Get(r, SessionName)
	if err != nil {
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
		return
	}
	addAuthContext(session, ctx)

	err = handler.ListTemplate.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, TemplateFailed(err), http.StatusInternalServerError)
	}
}

// EditRadio edits radios
func (handler *RadioHandler) EditRadio(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	name := query.Get("name")
	if name == "" {
		http.Error(w, "must specify a name", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		handler.ShowEditRadio(w, r, name)
	} else if r.Method == http.MethodPost {
		handler.DoEditRadio(w, r, name)
	}
}

// ShowEditRadio shows a form to edit a radio
func (handler *RadioHandler) ShowEditRadio(w http.ResponseWriter, r *http.Request, name string) {

	url := GetRadioURL(handler.Config.APIUrl, name)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, APIError(err), http.StatusInternalServerError)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, IOError(err), http.StatusInternalServerError)
		return
	}

	var radios []api.Radio
	err = json.Unmarshal(buf, &radios)
	if err != nil {
		http.Error(w, JSONError(err), http.StatusInternalServerError)
		return
	}

	ctx := pongo2.Context{}
	ctx["radio"] = radios[0]

	session, err := handler.Store.Get(r, SessionName)
	if err != nil {
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
		return
	}
	addAuthContext(session, ctx)

	err = handler.EditTemplate.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, TemplateFailed(err), http.StatusInternalServerError)
	}
}

// DoEditRadio handles a form submission
func (handler *RadioHandler) DoEditRadio(w http.ResponseWriter, r *http.Request, name string) {

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
		http.Error(w, JSONError(err), http.StatusInternalServerError)
		return
	}

	url := GetRadioURL(handler.Config.APIUrl, name)
	reader := bytes.NewReader(buf)
	req, _ := http.NewRequest("PATCH", url, reader)
	req.Header.Set("Authorization", handler.Config.APIToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, APIError(err), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, APIError(err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/radios/list", http.StatusTemporaryRedirect)
}
