package strutils

import (
	"encoding/json"
	"net/http"
	"strings"
)

var (
	threeWords = []string{"cok", "dik", "sht", "cnt", "kkk", "lsd", "cpd", "ira", "mmo", "ass", "fuc", "fuk", "fuq", "fux", "fck", "coc", "cok", "coq", "kox", "koc", "kok", "koq", "cac", "cak", "caq", "kac", "kak", "kaq", "dic", "dik", "diq", "dix", "dck", "pns", "psy", "fag", "fgt", "ngr", "nig", "cnt", "knt", "sht", "dsh", "twt", "bch", "cum", "clt", "kum", "klt", "suc", "suk", "suq", "sck", "lic", "lik", "liq", "lck", "jiz", "jzz", "gay", "gey", "gei", "gai", "vag", "vgn", "sjv", "fap", "prn", "jew", "joo", "gvr", "pus", "pis", "pss", "snm", "tit", "fku", "fcu", "fqu", "hor", "slt", "jap", "wop", "kik", "kyk", "kyc", "kyq", "dyk", "dyq", "dyc", "kkk", "jyz", "prk", "prc", "prq", "mic", "mik", "miq", "myc", "myk", "myq", "guc", "guk", "guq", "giz", "gzz", "sex", "sxx", "sxi", "sxe", "sxy", "xxx", "wac", "wak", "waq", "wck", "pot", "thc", "vaj", "vjn", "nut", "std", "lsd", "poo", "azn", "pcp", "dmn", "orl", "anl", "ans", "muf", "mff", "phk", "phc", "phq", "xtc", "tok", "toc", "toq", "mlf", "rac", "rak", "raq", "rck", "sac", "sak", "saq", "pms", "nad", "ndz", "nds", "wtf", "sol", "sob", "fob", "sfu"}
)

func IsBadWord(s string) bool {
	if len(s) == 3 {
		return isBadWord3Letters(s)
	}
	return isBadWord(s)
}

func isBadWord(s string) bool {
	resp, err := http.Get("https://www.purgomalum.com/service/json?text=" + s)
	if err != nil {
		return false
	}

	respBody := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return false
	}

	return strings.Contains(respBody["result"], "*")
}

func isBadWord3Letters(s string) bool {
	for _, si := range threeWords {
		if strings.Contains(s, si) {
			return true
		}
	}

	return false
}
