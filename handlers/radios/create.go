package radios

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/dabbotorg/radio-api/api"
	"github.com/flosch/pongo2"
)

// ViewCreate renders the template to create a radio
func (handler *Handler) ViewCreate(w http.ResponseWriter, r *http.Request) {
	ctx := pongo2.Context{}
	ctx["action"] = "create"

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	utils.AddAuthContext(session, ctx, handler.Config)
	utils.AddMetaContext(ctx, handler.Meta)

	radio := api.Radio{}
	radio.State = "ESCROW"
	/*data, err := json.Marshal(radio)
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}*/
	ctx["radio"] = radio

	err = handler.Templates.Edit.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, utils.TemplateFailed(err), http.StatusInternalServerError)
		return
	}
}

// Create handles a request to create a radio
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	raw, err := json.Marshal(radio)
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}
	buf := bytes.NewBuffer(raw)

	url := utils.CreateRadiosURL(handler.Config.APIUrl)
	req, _ := http.NewRequest(http.MethodPost, url, buf)
	req.Header.Set("Authorization", handler.Config.APIToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		http.Error(w, utils.APIError(err, resp), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/radios/list/valid", http.StatusTemporaryRedirect)
}
