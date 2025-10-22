package main

import (
	"context"
	"log"
	"os"
	"time"

	"governo/internal/auditlog"
	logplugin "governo/internal/plugins/log"
	webhookplugin "governo/internal/plugins/webhook"
	"governo/internal/rules"
	"governo/internal/server"
	"governo/internal/telemetry"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules []rules.Rule `yaml:"rules"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func watchConfig(path string, reload chan<- bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Printf("failed to close watcher: %v", err)
		}
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified config file:", event.Name)
					reload <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan struct{})
}

func main() {
	ctx := context.Background()
	cleanup := telemetry.Init(ctx)
	defer cleanup()

	cfg, err := loadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	// Register plugins
	rules.RegisterAction(&logplugin.LogAction{})
	rules.RegisterAction(&webhookplugin.WebhookAction{})

	dispatcher := auditlog.NewDispatcher(cfg.Rules)

	reload := make(chan bool)
	go watchConfig("config.yaml", reload)

	go func() {
		for {
			<-reload
			log.Println("Reloading config")
			time.Sleep(5 * time.Second) // add a small buffer to avoid reload loops
			cfg, err := loadConfig("config.yaml")
			if err != nil {
				log.Printf("error reloading config: %v", err)
				continue
			}
			dispatcher.UpdateRules(cfg.Rules)
		}
	}()

	srv := server.New(dispatcher)
	srv.Run()
}
