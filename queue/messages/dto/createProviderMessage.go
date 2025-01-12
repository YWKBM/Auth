package dto

type CreateProviderMessage struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	Region     string `json:"region"`
}
