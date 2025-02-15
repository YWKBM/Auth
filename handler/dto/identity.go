package dto

type IdentityRequest struct {
	Role      string `json:"role"`
	AuthToken string `json:"auth_token"`
}

type IdentityResponse struct {
	UserId int    `json:"userId"`
	Status string `json:"status"`
	Error  string `json:"error"`
}
