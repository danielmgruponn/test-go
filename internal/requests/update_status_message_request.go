package requests

type UpdateStatusMessage struct {
	Event			string `json:"event"`
	MessageId       int	   `json:"mns_id"`
}