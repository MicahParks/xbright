[![Go Report Card](https://goreportcard.com/badge/gitlab.com/MicahParks/xbright)](https://goreportcard.com/report/gitlab.com/MicahParks/xbright)
# xbright

A simple tool to edit your display's brightness.

Works for most Linux systems such as Debian based (Ubuntu, Mint), Red Hat based (CentOS, fedora), Arch, and others.
Any system were [`xrandr`](https://wiki.archlinux.org/index.php/Xrandr) works.

## Description

The `sliders` tab is where you'll find your displays as reported by `xrandr` followed by a slider and the current
brightness of the display. Use the slider to change the brightness of that display.

Get to the `settings` tab by clicking on it or pressing `Ctrl + Tab`. Using the appropriate radio button and preset
button, you can save or load your monitor's brightness settings from disk. If default is set, it will load every time
the application is started.

You can find settings saved to a simple file located at `~/.xbright.json`.

## Screenshots

![sliders tab](pics/sliders.png)

![settings tab](pics/settings.png)

Inspired by [https://github.com/LordAmit/Brightness](https://github.com/LordAmit/Brightness)
