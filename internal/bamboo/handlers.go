package bamboo

import (
	"encoding/json"
	"net/http"
)

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var answer = map[string]string{"status": "UP"}

	if err := json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) handleOrder(w http.ResponseWriter, r *http.Request) {
	type tmp struct {
		ICCIDs []string `json:"iccids"`
	}

	var temp tmp

	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		a.Logger.Errorf("failed to decoding request body: %s", newError(http.StatusInternalServerError, err))
		return
	}

	profiles, err := a.profileService.OrderProfile(r.Context(), temp.ICCIDs)
	if err != nil {
		return
	}

	Respond(w, http.StatusOK, profiles)
}
