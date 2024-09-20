package handlers

import (
	"github.com/kish1n/GiAuth/internal/data"
	"github.com/kish1n/GiAuth/internal/service/requests"
	"github.com/kish1n/GiAuth/internal/service/security"
	"github.com/kish1n/GiAuth/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"time"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRegistration(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	usernameOnce, err := UsersQ(r).FilterByUsername(req.Data.ID).Get()
	if err != nil {
		Log(r).WithError(err).Error("Error filter by username")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if usernameOnce != nil {
		Log(r).Errorf("User with this username: %s, already exit", req.Data.ID)
		ape.RenderErr(w, problems.Conflict())
		return
	}

	if !CheckAge(14, req.Data.Attributes.Birthday) {
		Log(r).Errorf("User younger than %v user's date of birth %s", 14, req.Data.Attributes.Birthday)
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	hashPassword, err := security.HashPassword(req.Data.Attributes.Password)
	if err != nil {
		Log(r).WithError(err).Error("Password hash error")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	user := data.User{
		Username:     req.Data.ID,
		Email:        req.Data.Attributes.Email,
		PasswordHash: hashPassword,
		FirstName:    req.Data.Attributes.FirstName,
		LastName:     req.Data.Attributes.LastName,
		MiddleName:   req.Data.Attributes.MiddleName,
		Birthday:     req.Data.Attributes.Birthday,
	}

	err = UsersQ(r).Insert(user)
	if err != nil {
		Log(r).WithError(err).Error("Error ger request NewDailyQuestion")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	newUser, err := UsersQ(r).FilterByUsername(req.Data.ID).Get()
	if err != nil {
		Log(r).WithError(err).Error("")
	}

	ape.Render(w, SuccessUserReg(newUser))
}

func CheckAge(age int, birthDate time.Time) bool {
	now := time.Now().UTC()
	ageDate := birthDate.AddDate(age, 0, 0)
	return !now.Before(ageDate)
}

func SuccessUserReg(user *data.User) resources.SuccessAuthResponse {
	return resources.SuccessAuthResponse{
		Data: resources.SuccessAuth{
			Key: resources.Key{
				ID:   user.Username,
				Type: resources.SUCCESS_REG,
			},
			Attributes: resources.SuccessAuthAttributes{
				Email:      user.Email,
				FirstName:  user.FirstName,
				LastName:   user.LastName,
				MiddleName: user.MiddleName,
				Username:   user.Username,
			},
		},
	}
}
