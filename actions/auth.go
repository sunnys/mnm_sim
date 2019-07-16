package actions

import (
	"database/sql"
	"strings"
	"encoding/json"

	"mnm_sim/models"
	"fmt"
	"time"
	"strconv"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
	"encoding/base64"
)

// AuthNew loads the signin page
func AuthNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(200, r.HTML("auth/new.html"))
}

// UsersCreate registers a new user with the application.
// AuthenticateUser godoc
// @Summary Add a User
// @Description add by json User
// @Tags users
// @Accept  json
// @Produce  json
// @Param User body models.AuthenticateUser true "Add User"
// @Success 200 {object} models.User
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /signin [post]
// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(u.Email))).First(u)
	// helper function to handle bad attempts
	bad := func() error {
		c.Set("user", u)
		verrs := validate.NewErrors()
		verrs.Add("email", "invalid email/password")
		c.Set("errors", verrs)
		return c.Render(422, r.Auto(c, u))
	}

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with the supplied email address.
			fmt.Printf("In-correct email.\n")
			return bad()
		}
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		fmt.Printf("Password mismatch.\n")
		return bad()
	}

	// --------------------------------------------------------------
	//  Block to generate tokens and append them to response header
	// --------------------------------------------------------------
	client, token, expiry, e := generate_tokens(32, 30)
	if e != nil {
		return bad()
	}
	tokens:= u.Tokens
	merged_token := make(map[string]string)
	token_hash, token_hash_err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if token_hash_err != nil {
		return errors.WithStack(token_hash_err)
	}
	merged_token[client] = string(token_hash)
	json.Unmarshal([]byte(tokens), &merged_token)
	json_merged_token, _:= json.Marshal(merged_token)
	u.Tokens = string(json_merged_token)
	update_err:= tx.Update(u)
	if update_err != nil {
		return errors.WithStack(update_err)
	}
	c.Response().Header().Set("uid", u.Email)
	c.Response().Header().Set("client", client)
	c.Response().Header().Set("access-token", token)
	c.Response().Header().Set("expiry", expiry)

	c.Flash().Add("success", "Welcome Back to Buffalo!")
	user := u.View()
	return c.Render(201, r.Auto(c, user))
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out!")
	return c.Redirect(302, "/")
}

// GenerateRandomBytes returns securely generated random bytes. 
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
    // Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generate_tokens(token_size int, expiry_time int) (string, string, string, error) {
	client, err := GenerateRandomString(token_size)
	token, err1 := GenerateRandomString(token_size)
	expiry := strconv.FormatInt((time.Now().AddDate(0,0, expiry_time).UnixNano() / 1000000), 10)
	if err != nil {
		return "", "", "", err
	}
	if err1 != nil {
		return "", "", "", err1
	}
	return client, token, expiry, nil
}