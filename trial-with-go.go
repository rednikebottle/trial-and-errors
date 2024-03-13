package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

type User struct {
	Login   string `json:"login"`
	HtmlUrl string `json:"html_url"`
}

var (
	name     string
	userData User
	option   string
)

func form() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's the username?").
				Value(&name).
				Validate(func(name string) error {
					res, err := http.Get("https://api.github.com/users/" + name)
					if err != nil {
						return err
					}

					if res.StatusCode == 404 {
						return errors.New("User not found")
					}

					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to see?").
				Options(
					huh.NewOption("About", "about"),
					huh.NewOption("Followers", "followers"),
					huh.NewOption("Following", "following"),
					huh.NewOption("Gists", "gists"),
					huh.NewOption("Starred repos", "starred"),
					huh.NewOption("Subscriptions", "subscriptions"),
					huh.NewOption("Organizations", "organizations"),
					huh.NewOption("Repos", "repos"),
					huh.NewOption("Events", "events"),
					huh.NewOption("Received Events", "received_events"),
					huh.NewOption("Quit", "quit"),
				).
				Validate(func(option string) error {
					if option == "quit" || option == "about" {
						return nil
					}

					var optionData []any
					getData("https://api.github.com/users/"+name+"/"+option, &optionData)

					if len(optionData) == 0 {
						return errors.New(option + ": Not found")
					}

					return nil
				}).
				Value(&option),
		),
	)

	err := form.Run()
	if err != nil {
		log.Error(err)
	}
}

func getData(url string, store any) {
	res, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(store)
}

func main() {
	logger := log.New(os.Stderr)
	form()
	getData("https://api.github.com/users/"+name, &userData)

	if option == "about" {
		logger.Infof(
			`Username: %s
GitHub URL: %s`, userData.Login, userData.HtmlUrl)
	}

	if option == "followers" {
		var followers []User
		getData("https://api.github.com/users/"+name+"/followers", &followers)

		for _, follower := range followers {
			logger.Infof(`%s - %s`, follower.Login, follower.HtmlUrl)
		}
	}
	if option == "following" {
		var following []User
		getData("https://api.github.com/users/"+name+"/following", &following)

		for _, followingUser := range following {
			logger.Infof(`%s - %s`, followingUser.Login, followingUser.HtmlUrl)
		}
	}

}
