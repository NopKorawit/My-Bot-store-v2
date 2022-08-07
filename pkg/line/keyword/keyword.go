package keyword

type Keyword string

const (
	Flavor = "Flavor"
)

func IsMenu(keyword string) bool {
	return keyword == "Relx INFINITY" || keyword == "INFY" || keyword == "JUES" || keyword == "INFY 7-11" || keyword == "INFINITE BOLD" || keyword == "All Flavor"
}
