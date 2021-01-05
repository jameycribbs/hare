package hare

type Contact struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func (c *Contact) GetID() int {
	return c.ID
}

func (c *Contact) SetID(id int) {
	c.ID = id
}

func (c *Contact) AfterFind() {
	*c = Contact(*c)
}
