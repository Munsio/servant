package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/rjz/githubhook.v0"
)

type server struct {
	Config     *Config
	Repository *Repository
}

type webhookHandler struct {
	Server *server
}

func (h *webhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	secret := []byte(h.Server.Config.WebhookSecret)
	hook, err := githubhook.Parse(secret, r)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}

	if hook.Event == "push" {
		h.Server.Repository.Pull()
	}

}

// NewServer creates the http listener
func NewServer(conf *Config, repo *Repository) (*server, error) {
	return &server{
		Config:     conf,
		Repository: repo,
	}, nil
}

func (serv *server) Run() {
	fs := http.FileServer(http.Dir(serv.Config.Path))

	if serv.Config.WebhookEnabled {
		hookHandler := &webhookHandler{
			Server: serv,
		}
		http.Handle("/webhook", hookHandler)
	}

	http.Handle("/", fs)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
