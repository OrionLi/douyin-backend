package model

type ActionRequest struct {
	Token      string `json:"token"`
	ToUserID   string `json:"to_user_id"`
	ActionType string `json:"action_type"`
	Content    string `json:"content"`
}
