package history

type History struct {
	Datetime int64  `json:"datetime"`
	From     string `json:"from"`
	To       string `json:"to"`
	Message  string `json:"message"`
}
