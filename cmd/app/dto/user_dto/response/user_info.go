package response

type UserInfo struct {
	UserID      int32  `json:"userID"`
	UserName    string `json:"userName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}
