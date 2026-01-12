package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Payload struct {
	Ref string `json:"ref"`
}

type App struct {
	Addr       string `mapstructure:"address"`
	Branch     string `mapstructure:"ref"`
	Script     string `mapstructure:"script"`
	InitScript string `mapstructure:"init_script"`
}

func (a *App) handle(w http.ResponseWriter, r *http.Request) {
	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Warnf("Failed to decode payload: %v", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	log.Infof("Received trigger for ref: %s", payload.Ref)

	shouldRun := payload.Ref == a.Branch
	defer func() {
		if shouldRun {
			go runScript(a.Script)
		}
	}()

	if shouldRun {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func runScript(path string) {
	cmd := exec.Command("/bin/sh", path)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	cmd.Run()

	log.Info(buf.String())

}

func main() {
	log.Info("Starting webhook server")

	var app App
	err := viper.Unmarshal(&app)
	if err != nil {
		log.Fatalf("Unable to decode config, %v", err)
	}

	log.Info("Running init script")
	runScript(app.InitScript)
	log.Info("Running first build")
	runScript(app.Script)

	log.Info("Server started, using config:")
	log.Infof("- Address=%s", app.Addr)
	log.Infof("- Branch=%s", app.Branch)
	log.Infof("- Script=%s", app.Script)

	http.HandleFunc("/webhook", app.handle)
	http.ListenAndServe(app.Addr, nil)
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "MST 2006-01-02 15:04:05",
	})

	// config file settings
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// default value
	viper.SetDefault("address", ":8001")
	viper.SetDefault("REF", "refs/heads/hugo-src")
	viper.SetDefault("script", "./scripts/build.sh")
	viper.SetDefault("init_script", "./scripts/init.sh")

	// env config
	viper.SetEnvPrefix("WEBHOOK")
	_ = viper.BindEnv("address", "ADDRESS")
	_ = viper.BindEnv("ref", "REF")
	_ = viper.BindEnv("script", "SCRIPT")
	_ = viper.BindEnv("init_script", "INIT_SCRIPT")

	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Warn("Unable to load .env file")
	}

	// read config file
	err = viper.ReadInConfig()
	if err != nil {
		log.Warn("Unable to load config file")
	}
}
