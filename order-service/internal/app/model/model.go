package model

// type Model struct {
// 	OrderCode   string       `json:"order_code,omitempty"`
// 	Date        string       `json:"date,omitempty"`
// 	ProductItem *ProductItem `json:"product_item"`
// }

type Model struct {
	ProductCode string `json:"id"`
	Name        string `json:"name"`
	Count       int    `json:"count"`
	Cost        int    `json:"cost"`
}
