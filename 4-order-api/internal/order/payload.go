package order

type OrderRequest struct {
	Products []string `json:"products"`
}
