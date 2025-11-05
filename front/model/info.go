package model

type Pedido struct {
	Usuario string `json:"usuario"`
	Total   int    `json:"total"`
	Estado  int    `json:"estado"`
}
