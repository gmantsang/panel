package radios

import (
	"net/http"
)

// Valid handles a confirmation to flag a radio as valid
func (handler *Handler) Valid(w http.ResponseWriter, r *http.Request) {
	handler.changeState(w, r, "VALID")
	http.Redirect(w, r, "/radios/list/escrow", http.StatusTemporaryRedirect)
}
