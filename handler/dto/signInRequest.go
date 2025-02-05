package dto

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *SignInRequest) Validate() bool {
	return s.Login != "" && s.Password != ""
}
