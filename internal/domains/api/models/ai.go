package models

type ModerateRequest struct {
	Text string `json:"text"`
}

type ProposeRequest struct {
	Advertiser string `json:"advertiser"`
	Title      string `json:"title"`
}
