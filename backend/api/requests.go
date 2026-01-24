package api

// JSON format from login body request
type LoginReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
	// called username on frontend to trick bots into filling out
	HoneyPot string `json:"username"`
}

// JSON format from login body request
type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
