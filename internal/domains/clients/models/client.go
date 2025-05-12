package models

type Gender string

const (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
)

type Client struct {
	ClientID string `json:"client_id"`
	Login    string `json:"login"`
	Age      int32  `json:"age"`
	Location string `json:"location"`
	Gender   Gender `json:"gender"`
}

type ClientUpsert struct {
	ClientID string `json:"client_id"`
	Login    string `json:"login"`
	Age      int32  `json:"age"`
	Location string `json:"location"`
	Gender   Gender `json:"gender"`
}
