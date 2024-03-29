package actions

import (
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"mnm_sim/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Phase)
// DB Table: Plural (phases)
// Resource: Plural (Phases)
// Path: Plural (/phases)
// View Template Folder: Plural (/templates/phases/)

// PhasesResource is the resource for the Phase model
type PhasesResource struct {
	buffalo.Resource
}

// List gets all Phases. This function is mapped to the path
// GET /phases

// List Phases godoc
// @Summary List phases
// @Description get all phases
// @Accept  json
// @Produce  json
// @Param access-token header string true "Access Token of successful authentication"
// @Param client header string true "Client Header of successful authentication"
// @Param expiry header string true "Expiry Header of successful authentication"
// @Param uid header string true "uid of user"
// @Param q query string false "name search by q"
// @Success 200 {array} models.Phase
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /phases [get]
func (v PhasesResource) List(c buffalo.Context) error {
	fmt.Printf("Header: %s\n", c.Request().Header.Get("uid"))
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	phases := &models.Phases{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Phases from the DB
	if err := q.All(phases); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Response().Header().Set("X-Pagination", q.Paginator.String())

	return c.Render(200, r.JSON(phases))
}

// Show gets the data for one Phase. This function is mapped to
// the path GET /phases/{phase_id}

// ShowPhase godoc
// @Summary Show a phase
// @Description get string by ID
// @Tags phases
// @Accept  json
// @Produce  json
// @Param id path int true "Phase ID"
// @Success 200 {object} models.Phase
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /phases/{id} [get]
func (v PhasesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Phase
	phase := &models.Phase{}

	// To find the Phase the parameter phase_id is used.
	if err := tx.Find(phase, c.Param("phase_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, phase))
}

// Create adds a Phase to the DB. This function is mapped to the
// path POST /phases

// AddPhase godoc
// @Summary Add a Phase
// @Description add by json Phase
// @Tags phases
// @Accept  json
// @Produce  json
// @Param access-token header string true "Access Token of successful authentication"
// @Param client header string true "Client Header of successful authentication"
// @Param expiry header string true "Expiry Header of successful authentication"
// @Param uid header string true "uid of user"
// @Param Phase body models.AddPhase true "Add Phase"
// @Success 200 {object} models.Phase
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /phases [post]
func (v PhasesResource) Create(c buffalo.Context) error {
	// Allocate an empty Phase
	phase := &models.Phase{}
	// fmt.Printf("Current User: %s", c.Get("current_user"))
	// Bind phase to the html form elements
	if err := c.Bind(phase); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(phase)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, phase))
	}

	// and redirect to the phases index page
	return c.Render(201, r.Auto(c, phase))
}

// Update changes a Phase in the DB. This function is mapped to

// the path PUT /phases/{phase_id}
// UpdatePhase godoc
// @Summary Update a phase
// @Description Update by json phase
// @Tags phases
// @Accept  json
// @Produce  json
// @Param  id path int true "Phase ID"
// @Param  phase body models.UpdatePhase true "Update phase"
// @Success 200 {object} models.Phase
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /phases/{id} [put]
func (v PhasesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Phase
	phase := &models.Phase{}

	if err := tx.Find(phase, c.Param("phase_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Phase to the html form elements
	if err := c.Bind(phase); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(phase)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, phase))
	}

	// and redirect to the phases index page
	return c.Render(200, r.Auto(c, phase))
}

// Destroy deletes a Phase from the DB. This function is mapped
// to the path DELETE /phases/{phase_id}
func (v PhasesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Phase
	phase := &models.Phase{}

	// To find the Phase the parameter phase_id is used.
	if err := tx.Find(phase, c.Param("phase_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(phase); err != nil {
		return errors.WithStack(err)
	}

	// Redirect to the phases index page
	return c.Render(200, r.Auto(c, phase))
}
