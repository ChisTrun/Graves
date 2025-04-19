package dto

import "time"

type ListOrders struct {
	From    *time.Time  `json:"from;omitempty"`
	To      *time.Time  `json:"to;omitempty"`
	OrderBy *SortMethod `json:"order_by;omitempty"`
}

type SortMethod struct {
	Column string `json:"column"`
	Order  string `json:"order"`
}
