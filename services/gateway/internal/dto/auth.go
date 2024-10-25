package dto

import "errors"

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r SignInRequest) Validate() error {
	if r.Login == "" {
		return errors.New("empty login or password")
	}

	if r.Password == "" {
		return errors.New("empty login or password")
	}

	return nil
}

type SignInResponse struct {
	Token string `json:"token"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}
