# Kodicast

This is a small [DIAL](http://www.dial-multiscreen.org) server that emulates
Chromecast-like devices, and implements the YouTube app. It proxies YouTube
commands from mobile app to Kodi YouTube plugin.

## Installation

I'm going to assume you're running Linux for this installation guide, preferably
Debian Jessie (or newer when their time comes). Debian before Jessie contains
too old versions of certain packages.

First, make sure you have the needed dependencies installed:

 *  golang 1.8+

These can be installed in one go under Debian Jessie (with jessie-backports):

    $ sudo apt-get install golang-1.8

If you haven't already set up a Go workspace, create one now. Some people like
to set it to their home directory, but you can also set it to a separate
directory. In any case, set the environment variable `$GOROOT` to this path:

    $ mkdir golang
    $ cd golang
    $ export GOPATH="`pwd`"

Then get the required packages and compile:

    $ go get -u github.com/sargo/kodicast

To run the server, run the executable `bin/kodicast` relative to your Go
workspace. Any Android phone with YouTube app (or possibly iPhone, but I haven't
tested) on the same network should recognize the server and it should be
possible to play the audio of videos on it. The Chrome extension doesn't yet
work.

    $ bin/kodicast

## Thanks

Big part of Kodicast is taken from
[Kodicast](https://github.com/sargo/kodicast) released on BSD license by
[Ayke van Laethem](https://sargo.nl/about).

I would like to thank the creators of
[leapcast](https://github.com/dz0ny/leapcast). Leapcast is a Chromecast
emulator, which was essential in the process of reverse-engineering the YouTube
protocol and better understanding the DIAL protocol.
