[![Go](https://github.com/nugget/autoslack/workflows/Go/badge.svg)](https://github.com/nugget/autoslack/actions?query=workflow%3AGo) [![Go Report Card](https://goreportcard.com/badge/github.com/nugget/autoslack)](https://goreportcard.com/report/github.com/nugget/autoslack)

# Nugget's Autoslack

This is a little service I run on my machine that updates my Slack status to
reflect whenever I am in a Highfive video call.  Published on principle, and
this should be easy enough to hack for your own use, but this really is just an
"itch scratching" tool.  I built it because I wanted it.  If it's useful for
you, even better.

PRs happily accepted.

## Installation

* `cp config.json.example config.json`
* Edit `config.json` to taste
* `make install`
* put this in your user crontab:
  ```
  @reboot $HOME/bin/autoslack_launch
  ```

The autolaunch script expects `tmux` to be installed.  That's not needed for
the service itself, but just how I chose to wrap it.
