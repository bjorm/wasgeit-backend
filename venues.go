package wasgeit

// Venue describes a place where Events take place
type Venue struct {
	ID        int64 `json:"-"`
	ShortName string
	Name      string
	URL       string
}
