# Garage Door Opener - hkgdo

Also known as hkgdo, or HomeKit Garage Door Opener.

## Table of Contents

- [Introduction](#introduction)
- [Technologies](#technologies)
- [Setup](#setup)
- [Running](#running)
- [Hardware](#hardware)
- [Todo](#todo)

## Introdution

This is yet another garage door opener project.
There are many options to look for on the web, I just thought about building my own solution for a RaspberryPi Zero that didn't need [HomeBridge][1] to work. The decision to not use HomeBridge is purely to get more familiar with HAP (HomeKit Accessory Protocol).

## Technologies

* Golang
* [periph.io - Go Hardware Library][2]
* [hc - Go HomeKit Accessory Protocol implementation][3]
* [fig - Go configuration management library][4]
* [Raspberry PI Zero W][6]

## Setup

Technically, all you need to build hkgdo is golang 1.11+ installed on your system. If you need help installing Go, visit [this page][5].

Also, you'll need `make` to be able to use some of the automation scripts. In case you don't have `make` available, just read the contents of the `Makefile`. Most commands are one-liners and very straight forward.

```shell
$ make
> Building release for ARM6 (Raspberry Pi Zero W)...
```
> **NOTE:** `make help` will give you a list of all targets available.

When build completes, you can find the binary in `./bin` folder.
That binary is an ARMv6 binary, so make sure you copy it to your Raspberry Pi first.

Some of the protocol and accessory settings that this project needs were added to a configuration file so there's no need to recompile in case of tweaking is necessary during runtime.
> **NOTE:** There is a sample configuration file in the root folder with lots of info on how to set it up. Name it `config.yaml` and place it next to the application's executable in Raspberry PI.

## Running

Once the binary and configuration files are copied over to the Raspberry Pi, navigate to the folder where the app is and run it as show below.

```shell
$ ./garagedoor-armv6
INFO 2021/02/22 22:53:49 ip_transport.go:184: Listening on port 57658
```

From an iPhone, pair with the device using the same PIN Number defined in your configuration file.
> **NOTE:** Make sure the Pi and the iPhone are connected to the same network.

## Hardware

As mentioned before, this was made to run on a [Raspberry Pi Zero W][6].
Other parts required to build this solution are:

### 4 Channel Relay Module

There are other solutions to provide relay functionality to your device that would work too. This is just the one I had around. No mistery here. All you have to do is to connect the HAT on your Raspberry PI Zero pins header and voila.

![4 Channel Relay Module Picture][7_relay_hat_pic]

Although the picture below shows a full size Raspberry Pi and the link provided for purchase also doesn't say anything about compatibility with the Pi Zero it's safe to use as both PIs share the same pinout.

A good and relialble source for finding out which pins do what is [this website][10].

> **NOTE:** Relay board on Amazon: [here][7]

### Reed switches

The reed switches are used to detect when the door is fully open or fully closed. My switches look a bit like the ones below, I assume these should do the trick too.

![Reed Switch Picture][8_reed_switch_pic]

For now, switches **MUST** be connected as follows:
* Fully Open Switch -> [GPIO_23](https://pinout.xyz/pinout/pin16_gpio23) + GND
* Fully Closed Sensor -> [GPIO_24](https://pinout.xyz/pinout/pin18_gpio24) + GND

Future releases could have this pins in configuration file.

> **NOTE:** Relay switches in Amazon: [here][8]

## Installation

### Software installation

Every Raspberry Pi Owner has it's own preferences when setting up their devices and I am not going to interfeer on that as the application has no dependencies other than enabling i2c modules on Raspberry Pi.

There are tons of websites and blogs that can teach you how to do that, so [here's my suggestion][i2c] on how to set it up.

### Hardware installation

Most tutorials and products out there teach you to jump wire your garage door motor by connecting some of the terminals on the back of the motor directly on the relay responsible for triggering the open/close events.

For compatibility issues that didn't work for me, so my solution was to use a spare garage door remote as trigger and connect the relay terminals to the terminals of the click button on the remote controller (there was some soldering going on here though...).

> **TODO**: Add picture here

It works exactly the same way, except that now my device doesn't have to be physically directly next to the motor, allowing it to sit in a conveinet location for maintenance, say, in case the sdcard goes bad fo wathever reason.

## Todo
- Actually use timeout when transitioning state
- Add Github Workflow for CI to build downloadable binaries
- Add tests
- Add GPIO settings to configuration file
- Add pictures of the actual installation
- Add suggested ansible setup and hardening playbooks
- Move Hardware details to a separate page to make main README cleaner

[1]: https://homebridge.io
[2]: https://periph.io
[3]: https://github.com/brutella/hc
[4]: https://githu.com/kkyr/fig
[5]: https://golang.org/doc/install
[6]: https://www.raspberrypi.org/products/raspberry-pi-zero-w/
[7]: https://www.amazon.com/MakerFocus-Raspberry-Expansion-Programming-Programmable/dp/B07Y54FKC6
[7_relay_hat_pic]: ./doc/relay_hat.jpg
[8]: https://www.amazon.com/Mbangde-Magnetic-Window-Contact-Personal/dp/B01M8JCGSO/ref=sr_1_9?dchild=1&keywords=reed+switches&qid=1614049890&s=electronics&sr=1-9
[8_reed_switch_pic]: ./doc/reed_switch.jpg
[10]: https://pinout.xyz
[i2c]: https://radiostud.io/howto-i2c-communication-rpi/
