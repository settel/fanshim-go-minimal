# fanshim-go-minimal

Absolutely bare driver for the Pimoroni Fan SHIM for Raspberry Pi 4, written in Go.

The compiled binary is only 2.1 MB in size and has no runtime dependencies (static binary). Easy to install and run as a background process, turning on and off the [Pimoroni Fan SHIM](https://shop.pimoroni.com/en-eu/products/fan-shim) as needed.

It draws ideas from [fanshim-go](https://github.com/jsiebens/fanshim-go) by jsiebens but is a minimal rewrite to run with Debian 13 (trixie).

# Installation

Either download and install the binary or compile on your own.

## Use pre-built binary

Download `fanshim-go-minimal` from the [GitHub Releases](https://github.com/settel/fanshim-go-minimal/releases)
page. Pushing a tag named `release-<version>` (for example, `release-1.0`)
builds the Linux ARM64 binary and attaches it to that release.

## Compile on your own

You can compile on your RaspPi 4 or cross-compile on your Linux PC.

### Prerequisites

- Go compiler (Debian/Ubuntu: install package golang-go)

### Compile

```
# automatically download and install compile-time dependencies
go mod tidy

# cross-compile for ARM64 architecture (RaspPi 4)
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o fanshim-go-minimal
```

# Setup

Install files to their locations and enable service in systemd.

```
sudo install -o root -g root -m 755 fanshim-go-minimal /usr/local/bin/
sudo install -o root -g root -m 644 fanshim-go-minimal.service /etc/systemd/system/

sudo systemctl daemon-reload
sudo systemctl enable --now fanshim-go-minimal
```
