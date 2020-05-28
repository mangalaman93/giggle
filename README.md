[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT) [![Build Status](https://travis-ci.com/mangalaman93/giggle.svg?branch=master)](https://travis-ci.com/mangalaman93/giggle) [![codecov](https://codecov.io/gh/mangalaman93/giggle/branch/master/graph/badge.svg)](https://codecov.io/gh/mangalaman93/giggle)

 [![Go Report Card](https://goreportcard.com/badge/github.com/mangalaman93/giggle)](https://goreportcard.com/report/github.com/mangalaman93/giggle) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/9933553bc3fb433d8d007cd917a64d90)](https://www.codacy.com/app/mangalaman93/giggle?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mangalaman93/giggle&amp;utm_campaign=Badge_Grade) [![Maintainability](https://api.codeclimate.com/v1/badges/b3e1b2f184edd8150ddd/maintainability)](https://codeclimate.com/github/mangalaman93/giggle/maintainability) [![CodeFactor](https://www.codefactor.io/repository/github/mangalaman93/giggle/badge)](https://www.codefactor.io/repository/github/mangalaman93/giggle)

## giggle

Sync overleaf repositories with GitHub (runs as a system tray)

## Requirements

* Only works for Linux and darwin OS (Mac)
* Doesn't work on windows yet, see https://github.com/getlantern/systray/issues/148

## Setup Config

* Copy the [example config file](https://github.com/mangalaman93/giggle/blob/master/config.json.example) to `$HOME/.config/.giggle/config.json` for Linux or `$HOME/Library/Application\ Support/.giggle/config.json` for Mac
* Add sync configuration to the `config.json` file
* **Make sure to create empty target repo on GitHub**

## Installation

### Linux

* Download the binary from [release](https://github.com/mangalaman93/giggle/releases/download/v0.1.0/giggle-linux-amd64) page
* Setup `config.json` as described in [setup config](#setup-config)
* `chmod +x giggle-linux-amd64`
* Execute the binary

## Mac

* Download the binary from [release](https://github.com/mangalaman93/giggle/releases/download/v0.1.0/giggle-darwin-amd64) page
* Setup `config.json` as described in [setup config](#setup-config)
* `chmod +x giggle-darwin-amd64`
* Execute the binary

## Building from Source

### Linux

```
sudo apt install libgtksourceview2.0-dev libgtk-3-dev libcairo2-dev libglib2.0-dev
go get -u github.com/mangalaman93/giggle
```

### Mac

```
xcode-select --install.
go get -u github.com/mangalaman93/giggle
```

## Future Work

* Add UI to manage configuration
* Fix running on windows
