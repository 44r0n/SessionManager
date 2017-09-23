package models

// ResponseData encapsulates a json response.
type ResponseData struct {
	Data Response `json:"Response"`
}

// Response to client
type Response struct {
	Status      int    `json:"Status"` //httpstatus
	Error       int    `json:"Error"`  //-1: unknown, -2: ecxeption, 1:ok, there are no 0's
	Description string `json:"Description"`
	Token       string `json:"Token"`
}
