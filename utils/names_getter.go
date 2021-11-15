package utils

import (
	"fmt"
	"math/rand"
)

type NamesGetter interface {
	GetNames() []string
	GetName() string
}

var names = []string{
	"Rustaceans",
	"The Mystery Krakens",
	"E = mc^2",
	"#include<team_name>",
	"Cherry Bombs",
	"dot",
	"sudo win",
	"untitled",
	"The Order of the Phoenix",
	"the real acm",
	"pika pika pikachu",
	"СССР",
	"ой мама",
	"The Cheems",
	"Lightning Round",
	"Rise of the Red",
	"Generals",
	"GIMMMME DA CUP!!",
	"Loading...",
	"{{ team.name }}",
	"new Team();",
	"Spray Painters",
	"Sales Xpress",
	"ÐOGE Hodlers",
	"The Worthless Wizards",
	"Fetching...",
	"Poodle Doodles",
	"g++ contest.o",
	"Nescafé",
	"redditors",
	"Exodius Minions",
	"tar -xpf team.gz",
	"Panzerkampf",
	"Bismarck",
	"Wind Ducks",
	"Thunderstorm",
	"Purple Cobras",
	"null",
	"Dream Team",
	"undefined",
}

type HardCodeNames struct {
	calledTimes int
	usedNames   map[int]bool
	numNames    int
}

func NewHardCodeNames() *HardCodeNames {
	return &HardCodeNames{
		calledTimes: 0,
		usedNames:   map[int]bool{},
		numNames:    len(names),
	}
}

func (h *HardCodeNames) GetNames() []string {
	return names
}

func (h *HardCodeNames) GetName() string {
getIndex:
	h.calledTimes++
	if h.calledTimes >= h.numNames {
		h.usedNames = map[int]bool{}
	}
	index := rand.Intn(h.numNames)

	if used, exists := h.usedNames[index]; !used || !exists {
		h.usedNames[index] = true
	} else {
		goto getIndex
	}

	return names[index]
}

type OrderedNames struct {
	lastIndex int
}

func NewOrderedNames() *OrderedNames {
	return &OrderedNames{-1}
}

func (o *OrderedNames) GetNames() []string {
	return nil
}

func (o *OrderedNames) GetName() string {
	o.lastIndex++
	return fmt.Sprintf("team%d", o.lastIndex)
}
