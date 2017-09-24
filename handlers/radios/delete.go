package radios

import (
	"net/http"
	"path"

	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/flosch/pongo2"
)

// ViewDelete renders a confirmation page to delete a radio
func (handler *Handler) ViewDelete(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)

	ctx := pongo2.Context{}
	ctx["name"] = name
	ctx["action"] = "delete"

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	utils.AddAuthContext(session, ctx, handler.Config)

	err = handler.Templates.Confirm.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, utils.TemplateFailed(err), http.StatusInternalServerError)
	}
}

// Delete handles a confirmation to delete a radio
func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)

	url := utils.GetRadioURL(handler.Config.APIUrl, name)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Authorization", handler.Config.APIToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, utils.APIError(err, resp), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/radios/list/valid", http.StatusTemporaryRedirect)
}
