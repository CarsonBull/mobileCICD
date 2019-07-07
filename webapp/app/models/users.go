package models

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	gormdb "github.com/revel/modules/orm/gorm/app"
	"golang.org/x/oauth2"
)

type User struct {
	UserId string `gorm:"primary_key"`
	Username string `gorm:"type:varchar(100);unique_index"`
	GithubAuth string
	Repos []Repo `gorm:"foreignkey:UId"`
}

type Repo struct {
	Id int64 `gorm:"primary_key"`
	UId string
	Reponame string
	Authorized bool
}


func CreateUser(username string) (User, error) {
	var user User

	// TODO: get a unique userid here. This can come from cognito
	user.UserId = username
	user.Username = username
	user.GithubAuth = ""
	user.Repos = nil

	gormdb.DB.Create(user)

	return user, nil
}

func GetUser(username string) User {
	var user User
	//fmt.Println(gormdb.DB.Table("users"))
	fmt.Println("Getting user")
	gormdb.DB.Where("username = ?", username).Preload("Repos").First(&user)

	return user
}

func (u *User) Authorize(token string) error {

	// TODO: check if its a valid token
	u.GithubAuth = token

	// TODO: get github repos here

	gormdb.DB.Save(u)

	return nil
}

func (u *User) Deauthorize() error {

	// TODO: check if its a valid token
	u.GithubAuth = ""

	for _, repo := range u.Repos {
		gormdb.DB.Delete(repo)
	}
	u.Repos = make([]Repo, 0)
	gormdb.DB.Save(u)

	return nil
}

func (u *User) SetupRepos() error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: u.GithubAuth},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)

	if err != nil {
		return err
	}


	for _, i := range repos {
		var repo Repo
		repo.Authorized = false
		repo.Reponame = i.GetName()
		repo.Id = i.GetID()
		u.Repos = append(u.Repos, repo)
	}

	gormdb.DB.Save(u)

	return nil
}

func (u *User) AuthorizeRepo(reponame string) error {

	if u.GithubAuth == "" {
		return errors.New("Github not authorized")
	}

	for i, _ := range u.Repos {
		if reponame == u.Repos[i].Reponame{
			u.Repos[i].Authorized = true
			gormdb.DB.Save(u)
			return nil
		}
	}

	return errors.New("Invalid repo")
}

func (u *User) CheckAuth() bool{

	return true
}