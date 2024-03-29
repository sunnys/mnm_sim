package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"mnm_sim/models"
	L "mnm_sim/lib"
	"fmt"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Build)
// DB Table: Plural (builds)
// Resource: Plural (Builds)
// Path: Plural (/builds)
// View Template Folder: Plural (/templates/builds/)

// BuildsResource is the resource for the Build model
type BuildsResource struct {
	buffalo.Resource
}

// List gets all Builds. This function is mapped to the path
// GET /builds

// List Builds godoc
// @Summary List builds
// @Description get all builds
// @Accept  json
// @Produce  json
// @Param access-token header string true "Access Token of successful authentication"
// @Param client header string true "Client Header of successful authentication"
// @Param expiry header string true "Expiry Header of successful authentication"
// @Param uid header string true "uid of user"
// @Param q query string false "name search by q"
// @Success 200 {array} models.Build
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /builds [get]
func (v BuildsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	builds := &models.Builds{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Builds from the DB
	if err := q.All(builds); err != nil {
		return errors.WithStack(err)
	}

	token := "e8ca8d3013172d0b956229eb8639f6e89b3aedf6f57038418fc68b04d27e1fe6"
	droplets, derr := L.ListDroplets(token, c)
	if derr != nil{
		return errors.WithStack(derr)
	}
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	// return c.Render(200, r.Auto(c, builds))
	return c.Render(200, r.Auto(c, droplets))
}

// Show gets the data for one Build. This function is mapped to
// the path GET /builds/{build_id}

// ShowBuild godoc
// @Summary Show a build
// @Description get string by ID
// @Tags builds
// @Accept  json
// @Produce  json
// @Param access-token header string true "Access Token of successful authentication"
// @Param client header string true "Client Header of successful authentication"
// @Param expiry header string true "Expiry Header of successful authentication"
// @Param uid header string true "uid of user"
// @Param id path string true "Build ID"
// @Success 200 {object} models.Build
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /builds/{id} [get]
func (v BuildsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Build
	build := &models.Build{}

	// To find the Build the parameter build_id is used.
	if err := tx.Find(build, c.Param("build_id")); err != nil {
		return c.Error(404, err)
	}
	token := "e8ca8d3013172d0b956229eb8639f6e89b3aedf6f57038418fc68b04d27e1fe6"
	droplet_id := 36136590
	droplet, derr := L.FindDroplet(token, droplet_id, c)
	if derr != nil{
		return errors.WithStack(derr)
	}
	// return c.Render(200, r.Auto(c, build))
	return c.Render(200, r.Auto(c, droplet))
}

// Create adds a Build to the DB. This function is mapped to the
// path POST /builds

// AddBuild godoc
// @Summary Add a Build
// @Description add by json Build
// @Tags builds
// @Accept  json
// @Produce  json
// @Param access-token header string true "Access Token of successful authentication"
// @Param client header string true "Client Header of successful authentication"
// @Param expiry header string true "Expiry Header of successful authentication"
// @Param uid header string true "uid of user"
// @Param Build body models.AddBuild true "Add Build"
// @Success 200 {object} models.Build
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /builds [post]
func (v BuildsResource) Create(c buffalo.Context) error {
	// Allocate an empty Build
	build := &models.Build{}

	// Bind build to the html form elements
	if err := c.Bind(build); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(build)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, build))
	}
	token := "e8ca8d3013172d0b956229eb8639f6e89b3aedf6f57038418fc68b04d27e1fe6"
	name := build.Name
	region := "blr1"
	size := "s-4vcpu-8gb"
	image := "kolearnDdServer3July2019"
	droplet, derr := L.CreateDroplet(token, name, region, size, image, c)
	if derr != nil{
		return errors.WithStack(derr)
	}
	fmt.Print(droplet)
	// and redirect to the builds index page
	return c.Render(201, r.Auto(c, build))
}

// Update changes a Build in the DB. This function is mapped to
// the path PUT /builds/{build_id}
func (v BuildsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Build
	build := &models.Build{}

	if err := tx.Find(build, c.Param("build_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Build to the html form elements
	if err := c.Bind(build); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(build)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, build))
	}

	// and redirect to the builds index page
	return c.Render(200, r.Auto(c, build))
}

// Destroy deletes a Build from the DB. This function is mapped
// to the path DELETE /builds/{build_id}
func (v BuildsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Build
	build := &models.Build{}

	// To find the Build the parameter build_id is used.
	if err := tx.Find(build, c.Param("build_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(build); err != nil {
		return errors.WithStack(err)
	}

	// Redirect to the builds index page
	return c.Render(200, r.Auto(c, build))
}
