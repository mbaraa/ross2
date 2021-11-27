package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type ContestAPI struct {
	endPoints map[string]http.HandlerFunc
	cRepo     data.ContestGetterRepo
}

func NewContestAPI(repo data.ContestGetterRepo) *ContestAPI {
	return (&ContestAPI{cRepo: repo}).initEndPoints()
}

func (c *ContestAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, trimContestNameSuffix(req.URL.Path), c.endPoints)
}

func trimContestNameSuffix(s string) string {
	path := strings.TrimPrefix(s, "/contest")
	if strings.Contains(path, "/single") {
		path = strings.TrimSuffix(path, path[len("/single/"):])
	}
	return path
}

func (c *ContestAPI) initEndPoints() *ContestAPI {
	c.endPoints = map[string]http.HandlerFunc{
		"GET /all/":    c.handleGetAllContests,
		"GET /single/": c.handleGetSingleContest,
	}
	return c
}

// GET /contest/all/
func (c *ContestAPI) handleGetAllContests(res http.ResponseWriter, req *http.Request) {
	contests, err := c.cRepo.GetAll()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(res).Encode(map[string]interface{}{
		"contests": contests,
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contest/single/{contestID}
func (c *ContestAPI) handleGetSingleContest(res http.ResponseWriter, req *http.Request) {
	contest, err := c.cRepo.Get(models.Contest{
		ID: getContestID(req.URL.Path),
	})
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if contest.TeamsHidden {
		contest.Teams = nil
	}

	err = json.NewEncoder(res).Encode(contest)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getContestID(path string) uint {
	id, _ := strconv.Atoi(path[len("/contest/single/"):])
	return uint(id)
}
