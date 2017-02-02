package models

// BoilerPlatedescribes a standard response from an API call. An additional
// "data" property exists which contains the data of a successful response but
// we do not define that here so that other structs can use this and add the
// typed data as needed rather than us using an interface{} here, it it the
// equivalent of:
// 	Data    interface{} `json:"data"`
type BoilerPlate struct {
	Context string   `json:"context"`
	Status  int      `json:"status"`
	Errors  []string `json:"error"`
}
