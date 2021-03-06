[![Go](https://github.com/nugget/autoslack/workflows/Go/badge.svg)](https://github.com/nugget/autoslack/actions?query=workflow%3AGo) [![Go Report Card](https://goreportcard.com/badge/github.com/nugget/autoslack)](https://goreportcard.com/report/github.com/nugget/autoslack)

# Nugget's Autoslack

This is a little service I run on my machine that updates my Slack status to
reflect whenever I am in a Highfive video call.  Published on principle, and
this should be easy enough to hack for your own use, but this really is just an
"itch scratching" tool.  I built it because I wanted it.  If it's useful for
you, even better.

PRs happily accepted.

### Configuration (Zoom)

* If you want to update your status whenever the Zoom client is running
  locally, you should trigger on the `zoom.us` process name.  If you only want
  your status updated when you are actively on a Zoom call, you should trigger
  on the `CptHost` process name.  Both process names are in the sample config
  file and you probably only want one of them.

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
