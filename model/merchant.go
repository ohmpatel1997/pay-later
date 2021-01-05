package model

type Merchant struct {
	Name     string
	Email    string
	Discount int //store as precision of 2 digits after decimal
}

func (m Merchant) TableName() string {
	return "merchant"
}

func (m Merchant) PrimaryKey() string {
	return m.Name
}

func (m Merchant) GetDiscountedAmount(amount int) float64 {
	return (float64(m.Discount) / float64(10000)) * float64(amount)
}
