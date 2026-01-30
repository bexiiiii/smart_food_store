package models

type Cart struct {
	ID     int
	UserID int
	Items  []CartItem
}

type CartItem struct {
	ProductID int
	Quantity  int
}
