# Nugget's Autoslack

This is a little service I run on my machine that updates my Slack status to
reflect whenever I am in a Highfive video call.

## Installation

* `cp config.json.example config.json`
* Edit `config.json` to taste
* `make install`
* put this in crontab:
  ```
  @reboot $HOME/bin/autoslack_launch
  ```
