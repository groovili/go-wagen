package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const keyENV = "APP_ENV"

func readConfig(env string) {
	if len(env) > 0 {
		env = fmt.Sprintf(".%s", env)
	}

	viper.SetConfigFile(fmt.Sprintf("./config/app%s.yml", env))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
}

func main() {
	app := "{{.Name}}"
	env := os.Getenv(keyENV)

	readConfig(env)

	log.Infof("Starting %s on %s env..", app, env)

	host, port := viper.GetString("host"), viper.GetString("port")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT)

	httpErr := make(chan error, 1)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	go func() {
		log.Infof("Started server on %s:%s..", host, port)
		httpErr <- http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
	}()

	select {
	case err := <-httpErr:
		log.Error(err)
	case <-stop:
		log.Info("Stopped via signal")
	}

	log.Infof("Stopping %s..", app)
}
