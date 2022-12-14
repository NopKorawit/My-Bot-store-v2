package keyword

import (
	"errors"
	"log"
)

type Keyword string

const (
	Flavor    = "Flavor"
	TypeNameA = "RELX INFINITY"
	TypeNameB = "INFY"
	TypeNameC = "JUES"
	TypeNameD = "INFY 7-11"
	TypeNameE = "INFINITE BOLD"
	TypeAll   = "ทั้งหมด"
	Story     = "คำขาย"
	Bank      = "บัญชี"
	Promptpay = "Promptpay"
	Help      = "Help"
)

var (
	ErrProductNotEnough = errors.New("product not enough")
	ErrCodenotFound     = errors.New("code not found")
)

func IsMenu(keyword string) bool {
	return keyword == TypeNameA || keyword == TypeNameB || keyword == TypeNameC || keyword == TypeNameD || keyword == TypeNameE || keyword == TypeAll
}

func ConvertType(keyword string) string {
	switch keyword {
	case TypeNameA:
		keyword = "A"
	case TypeNameB:
		keyword = "B"
	case TypeNameC:
		keyword = "C"
	case TypeNameD:
		keyword = "D"
	case TypeNameE:
		keyword = "E"
	default:
		log.Println("This Type not in Conditions")
	}
	return keyword
}
