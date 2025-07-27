package models

type UpdatePayload struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeletePayload struct {
	Name string `json:"name"`
}
