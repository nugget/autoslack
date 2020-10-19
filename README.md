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

The autolaunch script expects `tmux` to be installed.  That's not needed for
the service itself, but just how I chose to wrap it.
