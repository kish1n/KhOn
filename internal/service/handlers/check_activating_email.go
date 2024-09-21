package handlers

import (
	"github.com/go-chi/chi"
	"github.com/kish1n/GiAuth/internal/service/security"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strings"
)

func CheckActivatingEmail(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("user").(string)
	if !ok {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	code := strings.ToLower(chi.URLParam(r, "code"))

	if !security.CheckInEmailList(username, code) {
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	user, err := UsersQ(r).FilterByEmail(username).Get()
	if err != nil {
		Log(r).WithError(err).Error("Error get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if user == nil {
		Log(r).Errorf("User not found by id:%s", username)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	err = UsersQ(r).FilterByEmail(username).Update(map[string]any{
		"email_status": true,
	})
	if err != nil {
		Log(r).WithError(err).Error("Error update email status")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusFound)
}
