package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var (
	path           string
	repository     string
	command        string
	logLevel       string
	gitUser        string
	gitPass        string
	webhookEnabled bool
	webhookType    string
	webhookSecret  string
)

func init() {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	log.SetOutput(os.Stdout)

	flag.StringVar(&path, "path", "static", "Path where servant should serve files off")
	flag.StringVar(&repository, "repository", "", "The url of the repository where servant should pull its files")
	flag.StringVar(&command, "command", "", "An special command that should be executed after the pull is completed")
	flag.StringVar(&logLevel, "log-level", "info", "Verbosity of log output (debug,info,warn,error)")
	flag.StringVar(&gitUser, "git-user", "", "The git username (for private repos)")
	flag.StringVar(&gitPass, "git-pass", "", "The git password (for private repos)")
	flag.BoolVar(&webhookEnabled, "webhook-enabled", false, "Set to true if servant should update the repo on the webhook push command")
	flag.StringVar(&webhookType, "webhook-type", "github", "webhook provider (github)")
	flag.StringVar(&webhookSecret, "webhook-secret", "", "prevent webhook abuse")
	flag.Usage = printUsage
	flag.Parse()
}

func printUsage() {
	fmt.Println(`Usage: servant [options]
Options:`)
	flag.VisitAll(func(fg *flag.Flag) {
		fmt.Printf("\t--%s=%s\n\t\t%s\n", fg.Name, fg.DefValue, fg.Usage)
	})
}

func main() {

	if (len(os.Args) > 1) && (os.Args[1] == "help") {
		flag.Usage()
		os.Exit(1)
	}

	conf, err := initConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	repository, err := NewRepository(conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	server, err := NewServer(conf, repository)
	if err != nil {
		log.Fatal(err.Error())
	}

	server.Run()
}

/*
func initialize() {

	exists, err := direxists("./static/.git")

	if err != nil {
		log.Fatal(err)
	}

	if exists {
		pullRepo()
	} else {
		cloneRepo()
	}

	log.Printf("done")

}

func run() {
	fs := http.FileServer(http.Dir("static"))

	http.HandleFunc("/webhook", webhookParser)
	http.Handle("/", fs)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
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

func cloneRepo() {
	log.Printf("cloning repository")
	cmd := exec.Command("git", "clone", "https://github.com/Munsio/blog.git", "./static")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func pullRepo() {
	log.Printf("pulling repository")
	os.Chdir("./static")
	cmd := exec.Command("git", "pull")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir("../")
}

*/
