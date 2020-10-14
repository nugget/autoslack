package main

import (
	"encoding/json"
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

type AutoSlackConfig struct {
	SlackUserID  string `json:"slack_user_id"`
	SlackAPIKey  string `json:"slack_api_key"`
	LoopTime     int    `json:"loop_time"`
	Debug        bool   `json:"debug"`
	DefaultText  string `json:"default_status_text"`
	DefaultEmoji string `json:"default_status_emoji"`
}

type SlackStatus struct {
	Text  string `json:"status_text"`
	Emoji string `json:"status_emoji"`
}

var HIGHFIVE = SlackStatus{
	Text:  "Highfive call",
	Emoji: ":highfive:",
}

var DEFAULT = SlackStatus{
	Text:  "",
	Emoji: "",
}

func setStatus(api *slack.Client, newStatus SlackStatus) error {
	moo, err := api.GetUserProfile("U4XDDLMFY", true)
	if err != nil {
		return err
	}

	if newStatus.Text == moo.StatusText && newStatus.Emoji == moo.StatusEmoji {
		log.Debug("No change to status")
	} else {
		return api.SetUserCustomStatusWithUser(config.SlackUserID, newStatus.Text, newStatus.Emoji, 0)
	}

	return nil
}

func highFiveIsRunning() bool {
	plist, err := ps.Processes()
	if err != nil {
		log.WithError(err).Error("Can't get system process list")
		return false
	}

	for _, p := range plist {
		if p.Executable() == "Highfive" {
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

	DEFAULT.Text = c.DefaultText
	DEFAULT.Emoji = c.DefaultEmoji

	return c
}

func main() {
	initLog()
	config = loadConfig(path.Join(os.Getenv("HOME"), ".config", "autoslack", "config.json"))

	api := slack.New(config.SlackAPIKey)

	lastStatus := !highFiveIsRunning()

	for {
		highfiveStatus := highFiveIsRunning()
		if highfiveStatus == lastStatus {
			log.Debug("Highfive state is unchanged")
		} else {
			lastStatus = highfiveStatus

			setTo := DEFAULT

			if highfiveStatus {
				setTo = HIGHFIVE
			}

			err := setStatus(api, setTo)
			if err != nil {
				log.WithError(err).Error("Unable to set status")
			}
			log.WithFields(logrus.Fields{
				"text":  setTo.Text,
				"emoji": setTo.Emoji,
			}).Warning("Set Slack status")
		}
		time.Sleep(time.Duration(config.LoopTime) * time.Second)
		log.WithFields(logrus.Fields{
			"seconds": config.LoopTime,
		}).Trace("Looping")
	}
}
