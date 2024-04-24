package customer

import "fmt"

type Country struct {
	ID   int
	Name string
	Code string
}

func NewContry(id int, name string, code string) *Country {
	return &Country{
		ID:   id,
		Name: name,
		Code: code,
	}
}

func (c Country) String() string {
	return fmt.Sprintf("%s (%s)", c.Name, c.Code)
}

func (c Country) Equals(other Country) bool {
	return c.String() == other.String()
}
