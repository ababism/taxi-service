package location_client

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}