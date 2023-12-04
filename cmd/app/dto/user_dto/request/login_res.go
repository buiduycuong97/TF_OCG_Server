package request

type LoginRes struct {
	UserID       int32  `json:"userID"`
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	UserType     string `json:"userType"`
	Role         string `json:"role"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
