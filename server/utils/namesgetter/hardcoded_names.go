package namesgetter

import "math/rand"

type HardCodedNames struct {
	calledTimes int
	usedNames   map[int]bool
	names       []string
}

func NewHardCodedNames() *HardCodedNames {
	return &HardCodedNames{
		calledTimes: 0,
		usedNames:   map[int]bool{},
		names: []string{
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
		},
	}
}

func (h *HardCodedNames) GetNames() []string {
	return h.names
}

func (h *HardCodedNames) GetName() string {
getIndex:
	h.calledTimes++
	if h.calledTimes >= len(h.names) {
		h.calledTimes = 0
		h.usedNames = map[int]bool{}
	}
	index := rand.Intn(len(h.names))

	if used, exists := h.usedNames[index]; !used || !exists {
		h.usedNames[index] = true
	} else {
		goto getIndex
	}

	return h.names[index]
}
