package models

type UpdateMessage struct {
	Type    string
	Payload *Vehicle
	ID      *string
}
