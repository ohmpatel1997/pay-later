package model

type User struct {
	Name        string
	Email       string
	CreditLimit int //stored as cents, instead of dollars
	Dues        int
}

func (m User) TableName() string {
	return "user"
}

func (m User) PrimaryKey() string {
	return m.Name
}

func (m User) AllowAmount(amountToTransfer int) bool {

	if m.Dues+amountToTransfer > m.CreditLimit {
		return false
	}
	return true
}
