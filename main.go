package main

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"
	"path"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

var (
	log    *logrus.Logger
	config AutoSlackConfig
)

// AutoSlackConfig contains the configuration data loaded from the config file
type AutoSlackConfig struct {
	SlackUserID   string      `json:"slack_user_id"`
	SlackAPIKey   string      `json:"slack_api_key"`
	LoopTime      int         `json:"loop_time"`
	Debug         bool        `json:"debug"`
	DefaultStatus SlackStatus `json:"default_status"`
	States        []Trigger   `json:"states"`
}

// Trigger represents a process that, when seen locally, triggers an update to
// the users's Slack status
type Trigger struct {
	Process string      `json:"process"`
	Status  SlackStatus `json:"status"`
}

// SlackStatus embodies the text and icon associated with a specific Slack
// status update
type SlackStatus struct {
	Text  string `json:"text"`
	Emoji string `json:"emoji"`
}

func setStatus(api *slack.Client, newStatus SlackStatus) error {
	log.WithFields(logrus.Fields{
		"emoji": newStatus.Emoji,
		"text":  newStatus.Text,
	}).Debug("setStatus")

	moo, err := api.GetUserProfile(config.SlackUserID, true)
	if err != nil {
		return err
	}

	if newStatus.Text == moo.StatusText && newStatus.Emoji == moo.StatusEmoji {
		log.Debug("No change to status")
	} else {
		log.WithFields(logrus.Fields{
			"emoji": newStatus.Emoji,
			"text":  newStatus.Text,
		}).Info("Setting Slack Status")
		return api.SetUserCustomStatusWithUser(config.SlackUserID, newStatus.Text, newStatus.Emoji, 0)
	}

	return nil
}

func lookForProcess(name string) bool {
	log.WithFields(logrus.Fields{
		"name": name,
	}).Debug("lookForProcess")

	plist, err := ps.Processes()
	if err != nil {
		log.WithError(err).Error("Can't get system process list")
		return false
	}

	for _, p := range plist {
		if p.Executable() == name {
			return true
		}
	}

	return false
}

func initLog() {
	log = logrus.New()

	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_NOTICE, "autoslack")
	if err != nil {
		log.WithError(err).Error("Unable to connect to syslog")
	} else {
		log.AddHook(hook)
	}

}

func loadConfig(file string) (c AutoSlackConfig) {
	log.WithFields(logrus.Fields{
		"filename": file,
	}).Info("Loading configuration")
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.WithError(err).Fatal("Cannot load config file")
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&c)
	if err != nil {
		log.WithError(err).Fatal("Cannot parse config file")
	}

	if c.Debug {
		log.Info("Debug logging enabled")
		log.SetLevel(logrus.DebugLevel)
	}

	if c.Debug {
		fmt.Println("-- ")
		fmt.Printf("%+v\n", c)
		fmt.Println("-- ")
	}

	return c
}

func main() {
	var lastStatus string

	initLog()
	config = loadConfig(path.Join(os.Getenv("HOME"), ".config", "autoslack", "config.json"))

	api := slack.New(config.SlackAPIKey)

	for {
		found := false

		for _, trigger := range config.States {
			log.WithFields(logrus.Fields{
				"Process":    trigger.Process,
				"lastStatus": lastStatus,
			}).Debug("Looking for process")

			if lookForProcess(trigger.Process) {
				found = true

				if trigger.Process == lastStatus {
					log.Debug("Same Process")
				} else {
					log.Debug("New Process")

					err := setStatus(api, trigger.Status)
					if err != nil {
						log.WithError(err).Error("Unable to set status")
					} else {
						lastStatus = trigger.Process
					}
					break
				}
			} else {
				log.Debug("Process Not Found")
			}
		}

		if !found {
			log.Debug("Not Found")
			if lastStatus != "" {
				err := setStatus(api, config.DefaultStatus)
				if err != nil {
					log.WithError(err).Error("Unable to set status")
				} else {
					lastStatus = ""
				}
			}
		}

		log.Debug("-- ")
		time.Sleep(time.Duration(config.LoopTime) * time.Second)
		log.WithFields(logrus.Fields{
			"seconds": config.LoopTime,
		}).Trace("Looping")
	}
}
