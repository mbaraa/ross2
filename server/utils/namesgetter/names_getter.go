package namesgetter

import (
	"strings"
)

type NamesGetter interface {
	GetNames() []string
	GetName() string
}

func GetNamesGetter(nType string, names ...string) NamesGetter {
	switch strings.ToLower(nType) {
	case "ordered":
		return NewOrderedNames()
	case "random":
		return NewHardCodedNames()
	case "given":
		if names == nil {
			return nil
		}
		return NewGivenNames(names)
	}
	return NewHardCodedNames()
}
