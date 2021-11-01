package controllers

import (
	"github.com/mbaraa/ross2/config"
	"net/http"
)

func GetHandlerFromParentPrefix(res http.ResponseWriter, req *http.Request, endpointSuffix string, endpoints map[string]http.HandlerFunc) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", config.GetInstance().AllowedClients)
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if handler, exists := endpoints[req.Method+" "+endpointSuffix]; exists {
		handler(res, req)
		return
	}
	if req.Method != http.MethodOptions {
		http.NotFound(res, req)
	}
}
