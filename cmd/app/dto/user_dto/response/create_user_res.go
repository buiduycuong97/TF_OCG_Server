package response

type CreateUserRes struct {
	UserID   int32  `json:"userID"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	UserType string `json:"userType"`
}
