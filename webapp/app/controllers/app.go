package controllers

import (
	"fmt"
	"github.com/CarsonBull/mobileCICD/webapp/app/models"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {

	if username, err := c.Session.Get("username"); err == nil {
		return c.Render(username)
	}

	return c.Render()
}

// TODO: cleanup login.
func (c App) Login(username string) revel.Result{

	user := models.GetUser(username)

	fmt.Println(user)

	// TODO: remove this. This creates a user if one is not found.
	if user.Username == "" {
		revel.AppLog.Info("Error finding user")
		return c.Redirect("/create/" + username)
	}

	err := c.Session.Set("username", username)

	if err != nil {
		c.RenderError(err)
	}

	return c.Redirect("/")
}

func (c App) Logout() revel.Result{

	c.Session.Del("username")

	return c.Redirect("/")
}

func (c App) CreateAccount(username string) revel.Result {

	_, _ = models.CreateUser(username)

	return c.Redirect("/login/" + username)
}