package response

type CreateUserRes struct {
	UserID      int32  `json:"userID"`
	UserName    string `json:"userName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Role        string `json:"role"`
	UserType    string `json:"userType"`
}
