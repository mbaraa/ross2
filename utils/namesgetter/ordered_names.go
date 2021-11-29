package namesgetter

import "fmt"

type OrderedNames struct {
	lastIndex int
}

func NewOrderedNames() *OrderedNames {
	return &OrderedNames{0}
}

func (o *OrderedNames) GetNames() []string {
	return nil
}

func (o *OrderedNames) GetName() string {
	o.lastIndex++
	return fmt.Sprintf("team%d", o.lastIndex)
}
