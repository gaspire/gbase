package base

type SuccessResponse struct {
	Code       int         `json:"code"`
	Exp        int64       `json:"exp,omitempty"`
	Command    string      `json:"command,omitempty"`
	ServerTime string      `json:"server_time"`
	Tip        string      `json:"tip"`
	Output     interface{} `json:"output"`
}

type ErrResponse struct {
	Code       int    `json:"code"`
	Exp        int64  `json:"exp,omitempty"`
	Time       string `json:"time"`
	Tip        string `json:"tip"`
	ServerTime string `json:"server_time"`
}
