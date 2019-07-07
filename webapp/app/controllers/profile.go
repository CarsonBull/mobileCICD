package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CarsonBull/mobileCICD/webapp/app"
	"github.com/CarsonBull/mobileCICD/webapp/app/models"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
)

type Profile struct {
	*revel.Controller
}

func (c Profile) Profile() revel.Result {

	username, err := c.Session.Get("username")

	if err != nil {
	//	TODO redirect to a login or a create account page
		c.Redirect("/")
	}

	gitUrl := "https://github.com/login/oauth/authorize?client_id=" + app.GithubID +
		"&redirect_uri=" + app.CallbackUrl + "/profile/auth/callback/" + username.(string) + "&scope=repo"

	user := models.GetUser(username.(string))

	fmt.Println(user)

	return c.Render(user, username, gitUrl)
}

// TODO Handle errors
func (c Profile) Authorize(username string) revel.Result {

	uName, err := c.Session.Get("username")

	if err != nil {
		//	TODO redirect to a login or a create account page
		c.Redirect("/")
	} else if username != uName.(string) {
		c.Redirect("/logout")
	}

	user := models.GetUser(uName.(string))

	body := struct {
		ClientId string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code string `json:"code"`
	}{
		ClientId: app.GithubID,
		ClientSecret: app.GithubSecret,
		Code: c.Params.Query.Get("code"),
	}

	marshaledBody, _ := json.Marshal(body)

	reqs, _ := http.NewRequest("POST","https://github.com/login/oauth/access_token" , bytes.NewBuffer(marshaledBody))

	reqs.Header.Set("Accept", "application/json")
	reqs.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(reqs)

	if err != nil {
		fmt.Println(err)
		c.RenderError(err)
	}

	defer resp.Body.Close()
	code , _ := ioutil.ReadAll(resp.Body)

	var dat map[string]interface{}

	_ = json.Unmarshal(code, &dat)

	token := dat["access_token"]

	_ = user.Authorize(token.(string))

	_ = user.SetupRepos()

	return c.Redirect("/profile/" + uName.(string))

}


func (c Profile) AuthorizeRepo(reponame string) revel.Result {

	username, err := c.Session.Get("username")

	if err != nil {
		//	TODO redirect to a login or a create account page
		c.Redirect("/")
	}

	user := models.GetUser(username.(string))

	_ = user.AuthorizeRepo(reponame)

	return c.Redirect("/profile/" + username.(string))

}
