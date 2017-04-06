package main

import (
	log "github.com/Sirupsen/logrus"
	"net/url"
	"os"
	"os/exec"
)

// Repository represents an git repository
type Repository struct {
	Config *Config
	Url    *url.URL
	User   string
	Pass   string
}

// NewServer creates the http listener
func NewRepository(conf *Config) (*Repository, error) {

	url, err := url.Parse(conf.Repository)
	if err != nil {
		log.Fatal(err)
	}

	repo := Repository{
		Config: conf,
		Url:    url,
	}

	if len(conf.Repository) > 0 {
		exists, err := direxists(conf.Path + "/.git")
		if err != nil {
			log.Fatal(err)
		}

		if exists {
			repo.Pull()
		} else {
			repo.Clone()
		}
	}

	return &repo, nil
}

func (repo *Repository) Clone() {
	log.Printf("cloning repository")
	cmd := exec.Command("git", "clone", repo.Url.String(), repo.Config.Path)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (repo *Repository) Pull() {
	log.Printf("pulling repository")
	os.Chdir(repo.Config.Path)
	cmd := exec.Command("git", "pull")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir("../")
}

func direxists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
