package messengerbot

type rawError struct {
	Error Error `json:"error"`
}

type Error struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	ErrorData string `json:"error_data"`
	FbTraceId string `json:"fbtrace_id"`
}
