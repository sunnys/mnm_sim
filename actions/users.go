package actions

import (
	"mnm_sim/models"
	"fmt"
	"time"
	"strconv"
	"encoding/json"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
	fmt.Printf("1\n")
	return func(c buffalo.Context) error {
		if uid := c.Request().Header.Get("uid"); uid != "" {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Where("email = ?", uid).First(u)
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
	fmt.Printf("2\n")
	return func(c buffalo.Context) error {
		uid := c.Request().Header.Get("uid")
		client := c.Request().Header.Get("client")
		token := c.Request().Header.Get("access-token")
		expiry := c.Request().Header.Get("expiry")
		if uid == "" {
			c.Session().Set("redirectURL", c.Request().URL.String())
			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			return c.Render(401, r.JSON(map[string]string{"message": "You are not authorized to see that page!, please authenticate and provide correct headers in the request "}))
		}
		u := &models.User{}
		tx := c.Value("tx").(*pop.Connection)
		err := tx.Where("email = ?", uid).First(u)
		if err != nil {
			return c.Render(401, r.JSON(map[string]string{"message": "You are not authorized to see that page!, please authenticate and provide correct headers in the request "}))
		}
		token_matches:= MatchToken(u, client, token, expiry)
		if token_matches != true {
			return c.Render(401, r.JSON(map[string]string{"message": "You are not authorized to see that page!, please authenticate and provide correct headers in the request "}))
		}
		c.Response().Header().Set("uid", u.Email)
		c.Response().Header().Set("client", client)
		c.Response().Header().Set("access-token", token)
		c.Response().Header().Set("expiry", expiry)
		return next(c)
	}
}

func MatchToken(u *models.User, client, token, expiry string) bool {
	current_time := time.Now().UnixNano() / 1000000
	expiry_time, c_err := strconv.ParseInt(expiry, 10, 64)
	if c_err != nil || current_time > expiry_time {
		return false
	}
	tokens:= u.Tokens
	merged_token := make(map[string]string)
	json.Unmarshal([]byte(tokens), &merged_token)
	token_hash:= merged_token[client]
	err := bcrypt.CompareHashAndPassword([]byte(token_hash), []byte(token))
	if err != nil {
		return false
	}
	return true

}
