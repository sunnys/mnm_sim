package actions

import (
	"mnm_sim/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// Create new user godoc
// @Summary Create user
// @Description Create user form
// @Accept  json
// @Produce  json
// @Success 200 {object} web.SuccessObj
// @Router /users/new [get]
func UsersNew(c buffalo.Context) error {
	u := models.User{}
	c.Set("user", u)
	return c.Render(200, r.JSON(map[string]string{"message": "Welcome to Buffalo!"}))
}

// UsersCreate registers a new user with the application.
// AddUser godoc
// @Summary Add a User
// @Description add by json User
// @Tags users
// @Accept  json
// @Produce  json
// @Param User body models.AddUser true "Add User"
// @Success 200 {object} models.User
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /users [post]
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(200,  r.Auto(c, verrs))
	}

	c.Response().Header().Set("uid", u.Email)
	c.Flash().Add("success", "Welcome to Buffalo!")

	return c.Render(201, r.Auto(c, u))
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}
