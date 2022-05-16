package response

type Error struct {
	Msg string `json:"errorMsg"`
}

type Info struct {
	Msg string `json:"infoMsg"`
}