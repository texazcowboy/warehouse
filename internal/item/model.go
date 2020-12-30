package item

type Item struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"max=6"`
}
