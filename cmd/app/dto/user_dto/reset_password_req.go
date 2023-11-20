package user_dto

type ResetPasswordReq struct {
	ResetToken      string `json:"resetToken"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}
