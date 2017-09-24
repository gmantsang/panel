package radios

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/dabbotorg/radio-api/api"
	"github.com/flosch/pongo2"
)

// ViewList renders the radio-list template
func (handler *Handler) ViewList(w http.ResponseWriter, r *http.Request) {
	state := strings.ToUpper(path.Base(r.URL.Path))
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

	url := utils.ListRadiosURL(handler.Config.APIUrl, limit, offset, state)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, utils.APIError(err), http.StatusInternalServerError)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, utils.IOError(err), http.StatusInternalServerError)
		return
	}

	var radios api.RadioList
	err = json.Unmarshal(buf, &radios)
	if err != nil {
		http.Error(w, utils.JSONError(err), http.StatusInternalServerError)
		return
	}

	ctx := pongo2.Context{}
	ctx["radios"] = radios.Radios
	ctx["count"] = len(radios.Radios)
	ctx["total"] = radios.Count
	ctx["page"] = int32(math.Floor(float64(offsetInt+100) / float64(100)))
	ctx["pages"] = int32(math.Floor(float64(radios.Count+100) / float64(100)))
	ctx["state"] = state

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	utils.AddAuthContext(session, ctx, handler.Config)

	err = handler.Templates.List.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, utils.TemplateFailed(err), http.StatusInternalServerError)
	}
}
