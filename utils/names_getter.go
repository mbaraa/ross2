package utils

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
	lastIndex int
}

func NewHardCodeNames() *HardCodeNames {
	return &HardCodeNames{-1}
}

func (h *HardCodeNames) GetNames() []string {
	return names
}

func (h *HardCodeNames) GetName() string {
	h.lastIndex++
	if h.lastIndex >= len(names) {
		h.lastIndex = 0
	}
	return names[h.lastIndex]
}
