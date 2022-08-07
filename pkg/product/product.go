package product

type Product struct {
	Code     string
	Type     string
	Name     string
	Quantity int
}

func (p Product) GetQtySymbol() string {
	if p.Quantity < 0 {
		return "ระบบมีปัญหา"
	} else if p.Quantity == 0 {
		return "❌"
	} else if p.Quantity < 3 {
		return "⚠️"
	} else {
		return "✅"
	}
}
