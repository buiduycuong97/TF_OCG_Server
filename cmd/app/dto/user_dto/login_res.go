package user_dto

type LoginRes struct {
	UserID       int32  `json:"userID"`
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	UserType     string `json:"userType"`
	Role         string `json:"role"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
