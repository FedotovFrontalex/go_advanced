package stat

type StatGetResponse struct {
	Period string `json:"period"`
	Sum    int    `json:"sum"`
}
