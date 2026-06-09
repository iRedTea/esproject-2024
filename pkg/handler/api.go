package handler

import (
	"esproject"
	"esproject/pkg/util"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type projectRequest struct {
	Name    string   `json:"name"`
	OwnerId int64    `json:"owner_id"`
	Tags    []string `json:"tags"`
}

// project       godoc
// @Summary      Create new project
// @Tags         project
// @Produce      json
// @Param        project  body   handler.projectRequest  true       "Project data in JSON"
// @Param        Authorization   header    string        true   	"JWT Bearer token (authorization)"
// @Success      200 {object} esproject.Project "created project"
// @Router       /api/project [post]
func (h *Handler) create(c *gin.Context) {
	current, _ := c.Get(userCtx)
	var input projectRequest

	if (input.OwnerId != current) && (!util.MatchRole(c, "admin")) {
		util.NewErrorResponse(c, http.StatusForbidden, "no perms")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	project, err := h.repo.CreateProject(input.OwnerId, input.Name, input.Tags)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, project)
}

// project/all      godoc
// @Summary      Get ids of projects
// @Tags         project
// @Produce      json
// @Param        limit  query      int  true  "limit of projects size"
// @Param        start  query      int  true  "select offset"
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200 {array} int
// @Router       /api/project/all [get]
func (h *Handler) all(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	start, err := strconv.Atoi(c.Query("start"))
	if err != nil {
		start = 0
	}

	c.JSON(http.StatusOK, h.repo.GetProjects(start, limit))
}

type countResponse struct {
	Count int64 `json:"count"`
}

// user/count    godoc
// @Summary      Get count of projects
// @Tags         project
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200   {object}  handler.countResponse
// @Router       /api/project/count [get]
func (h *Handler) count(c *gin.Context) {
	c.JSON(http.StatusOK, countResponse{Count: h.repo.GetProjectsCount()})
}

// project/delete         godoc
// @Summary               Delete project by id (only for owner of project)
// @Tags                  project
// @Produce               json
// @Param                 id  path      int  true  "id of project"
// @Param                 Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success               200
// @Router                /api/project/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	project, err := h.repo.GetProject(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || ((project.Id != currentUser.Id) && (currentUser.Access != "admin")) {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	err = h.repo.DeleteProject(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "{ \"status\": \"ok\" }")
}

// project       godoc
// @Summary      Get project by id
// @Tags         project
// @Produce      json
// @Param        id  path      int  true  "id of project"
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200
// @Router       /api/project/{id} [get]
func (h *Handler) getProject(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	project, err := h.repo.GetProject(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, project)
}

// project       godoc
// @Summary      Edit project (only for owner of project)
// @Tags         project
// @Produce      json
// @Param        Authorization    header    string                           true   	"JWT Bearer token (authorization)"
// @Param        project     body      esproject.Project    true   	"edited project"
// @Success      200   {object} esproject.Project
// @Router       /api/project/ [post]
func (h *Handler) editProject(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	project, err := h.repo.GetProject(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || ((project.Id != currentUser.Id) && (currentUser.Access != "admin")) {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	var input esproject.Project
	if err := c.BindJSON(&input); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.Id = project.Id

	response, err := h.repo.SaveProject(&input)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, response)
}

type memberRequest struct {
	Project int64 `json:"project"`
	User    int64 `json:"user"`
}

// member       godoc
// @Summary      invite member (only for owner of project)
// @Tags         member
// @Produce      json
// @Param        Authorization    header    string                           true   	"JWT Bearer token (authorization)"
// @Param        project     body      handler.memberRequest    true   	"invite data in json format"
// @Success      200   {object}  esproject.ProjectInvite
// @Router       /api/member/invite/ [post]
func (h *Handler) invite(c *gin.Context) {
	var input memberRequest
	if err := c.BindJSON(&input); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	project, err := h.repo.GetProject(input.Project)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || ((input.User != project.OwnerId) && (currentUser.Access != "admin")) {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	invite, err := h.repo.CreateInvite(input.Project, input.User)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, invite)
}

// member       godoc
// @Summary      invites of current user
// @Tags         member
// @Produce      json
// @Param        Authorization    header    string  true   	"JWT Bearer token (authorization)"
// @Success      200 {array} int
// @Router       /api/member/invite/my [get]
func (h *Handler) myInvites(c *gin.Context) {

	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	invites, err := h.repo.FindInvites("user = ?", fmt.Sprintf("%d", currentUser.Id))
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var ans []int64
	for _, inv := range invites {
		ans = append(ans, inv.Id)
	}
	c.JSON(http.StatusOK, ans)
}

// member       godoc
// @Summary      accept invite
// @Tags         member
// @Produce      json
// @Param        Authorization    header    string  true   	"JWT Bearer token (authorization)"
// @Param        id  path      int  true  "id of invite"
// @Success      200
// @Router       /api/member/invite/accept/{id} [post]
func (h *Handler) accept(c *gin.Context) {
	current, _ := c.Get(userCtx)

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	invite, err := h.repo.GetInvite(id)
	if (err != nil) || (invite.User != current) {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	project, err := h.repo.GetProject(invite.Project)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	project.Members = append(project.Members, invite.User)

	err = h.repo.DeleteInvite(invite.Id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = h.repo.SaveProject(project)
	if err != nil {
		util.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "{ \"status\": \"ok\" }")
}

// member       godoc
// @Summary      kick member (only for owner of project)
// @Tags         member
// @Produce      json
// @Param        Authorization    header    string                           true   	"JWT Bearer token (authorization)"
// @Param        project     body      handler.memberRequest    true   	"invite data in json format"
// @Success      200
// @Router       /api/member/kick/ [post]
func (h *Handler) kick(c *gin.Context) {
	var input memberRequest
	if err := c.BindJSON(&input); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	project, err := h.repo.GetProject(input.Project)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || ((input.User != project.OwnerId) && (currentUser.Access != "admin")) {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	if !slices.Contains(project.Members, input.User) {
		util.NewErrorResponse(c, http.StatusBadRequest, "user is not in project")
		return
	}

	project.Members = removeElement(project.Members, input.User)
	c.JSON(http.StatusOK, "{ \"status\": \"ok\" }")
}

// project/my      godoc
// @Summary      Get ids of projects when user is member
// @Tags         project
// @Produce      json
// @Param        limit  query      int  true  "limit of projects size"
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200 {array} int
// @Router       /api/project/my [get]
func (h *Handler) my(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	current, _ := c.Get(userCtx)
	if err != nil {
		limit = 10
	}
	var result []int64
	for _, p := range h.repo.GetProjects(0, int(h.repo.GetProjectsCount())) {
		if (len(result) == limit) && (limit > 0) {
			break
		}
		project, err := h.repo.GetProject(p)
		if err != nil {
			util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if (project.OwnerId == current.(int64)) || containsElement(project.Members, current.(int64)) {
			result = append(result, project.Id)
		}
	}

	c.JSON(http.StatusOK, result)
}

// member       godoc
// @Summary      invite by id
// @Tags         member
// @Produce      json
// @Param        Authorization    header    string  true   	"JWT Bearer token (authorization)"
// @Success      200 {array} int
// @Router       /api/member/invite/{id} [get]
func (h *Handler) getInvite(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	invite, err := h.repo.GetInvite(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, invite)
}

func removeElement(slice []int64, value int64) []int64 {
	// Find the index of the value
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func containsElement(slice []int64, value int64) bool {
	// Find the index of the value
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
