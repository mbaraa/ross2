package namesgetter

import "math/rand"

type GivenNames struct {
	calledTimes int
	usedNames   map[int]bool
	names       []string
}

func NewGivenNames(names []string) *GivenNames {
	return &GivenNames{
		calledTimes: 0,
		usedNames:   map[int]bool{},
		names:       names,
	}
}

func (g *GivenNames) GetNames() []string {
	return g.names
}

func (g *GivenNames) GetName() string {
getIndex:
	g.calledTimes++
	if g.calledTimes >= len(g.names) {
		g.calledTimes = 0
		g.usedNames = map[int]bool{}
	}
	index := rand.Intn(len(g.names))

	if used, exists := g.usedNames[index]; !used || !exists {
		g.usedNames[index] = true
	} else {
		goto getIndex
	}

	return g.names[index]
}
