package attribute

type Assignment struct {
	Attribute   *Attribute
	AttributeID string
}

type BaseAssignedAttribute struct {
	Assignment *Assignment
}

func (b *BaseAssignedAttribute) Attribute() *Attribute {
	return b.Assignment.Attribute
}

func (b *BaseAssignedAttribute) AttributePk() string {
	return b.Assignment.AttributeID
}
