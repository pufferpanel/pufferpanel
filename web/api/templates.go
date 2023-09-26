package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
)

func registerTemplates(g *gin.RouterGroup) {
	g.Handle("GET", "", middleware.RequiresPermission(pufferpanel.ScopeTemplatesView), getRepos)
	g.Handle("POST", "", middleware.RequiresPermission(pufferpanel.ScopeTemplatesRepoCreate), addRepo)
	g.Handle("OPTIONS", "", response.CreateOptions("GET", "POST"))

	g.Handle("GET", "/:repo", middleware.RequiresPermission(pufferpanel.ScopeTemplatesView), getsTemplatesForRepo)
	g.Handle("DELETE", "/:repo", middleware.RequiresPermission(pufferpanel.ScopeTemplatesRepoDelete), deleteRepo)
	g.Handle("OPTIONS", "/:repo", response.CreateOptions("GET", "PUT", "DELETE"))

	g.Handle("GET", "/:repo/:name", middleware.RequiresPermission(pufferpanel.ScopeTemplatesView), getTemplateFromRepo)
	g.Handle("DELETE", "/0/:name", middleware.RequiresPermission(pufferpanel.ScopeTemplatesLocalEdit), deleteTemplate)
	g.Handle("PUT", "/0/:name", middleware.RequiresPermission(pufferpanel.ScopeTemplatesLocalEdit), putTemplate)
	g.Handle("OPTIONS", "/:repo/:name", response.CreateOptions("GET"))
	g.Handle("OPTIONS", "/0/:name", response.CreateOptions("GET", "DELETE", "PUT"))
}

// @Summary Get all repos
// @Description Gets all repos that are available to pull template from
// @Success 200 {object} []models.TemplateRepo
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/templates [get]
// @Security OAuth2Application[templates.view]
func getRepos(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	repos, err := ts.GetRepos()
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, repos)
}

// @Summary Get all templates from repo
// @Description Gets all templates from a repository
// @Param repo path uint true "Repo id"
// @Success 200 {object} []models.Template
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/templates/{repo} [get]
// @Security OAuth2Application[templates.view]
func getsTemplatesForRepo(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	repoId, err := cast.ToUintE(c.Param("repo"))
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	templates, err := ts.GetAllFromRepo(repoId)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, templates)
}

// @Summary Add repo
// @Description Adds a new repo to the service
// @Param repo body models.TemplateRepo true "Repo information"
// @Success 200 {object} models.TemplateRepo
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/templates [post]
// @Security OAuth2Application[templates.repo.create]
func addRepo(c *gin.Context) {
	var repo *models.TemplateRepo
	err := c.BindJSON(repo)

	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	if repo.Name == "" {
		response.HandleError(c, pufferpanel.ErrFieldRequired("repoName"), http.StatusBadRequest)
		return
	}

	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	err = ts.AddRepo(repo)

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	c.JSON(http.StatusOK, repo)
}

// @Summary Delete repo
// @Description Deletes a repo from the service
// @Param repo path uint true "Repo Id"
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/templates/{repo} [delete]
// @Security OAuth2Application[templates.repo.delete]
func deleteRepo(c *gin.Context) {
	repoId, err := cast.ToUintE(c.Param("repo"))
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	err = ts.DeleteRepo(repoId)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get template
// @Description Gets a template from the repo
// @Param repo path uint true "Repo Id"
// @Param template path string true "Template name"
// @Success 200 {object} models.Template
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Router /api/templates/{repo}/{template} [get]
// @Security OAuth2Application[templates.view]
func getTemplateFromRepo(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	repoId, err := cast.ToUintE(c.Param("repo"))
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	template, err := ts.Get(repoId, c.Param("name"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(404)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, template)
}

// @Summary Adds or updates a template
// @Success 204 {object} nil
// @Failure 400 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param template body pufferpanel.Server true "Template"
// @Param name path string true "Template name"
// @Router /api/templates/local/{name} [put]
// @Security OAuth2Application[templates.local.edit]
func putTemplate(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	templateName := c.Param("name")
	templateRequest := pufferpanel.Server{}
	err := c.MustBindWith(&templateRequest, binding.JSON)
	if response.HandleError(c, err, http.StatusBadRequest) {
		return
	}

	template, err := ts.Get(ts.GetLocalRepoId(), templateName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		template = &models.Template{
			Name: templateName,
		}
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	template.Server = templateRequest
	err = ts.Save(template)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Deletes template
// @Description Deletes template
// @Success 204 {object} nil
// @Failure 404 {object} pufferpanel.ErrorResponse
// @Failure 500 {object} pufferpanel.ErrorResponse
// @Param name path string true "Template name"
// @Router /api/templates/local/{name} [delete]
// @Security OAuth2Application[templates.repo.delete]
func deleteTemplate(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ts := &services.Template{DB: db}

	err := ts.Delete(c.Param("name"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}
