# AuralPi

Minimal Buildroot firmware for Raspberry Pi 4 as a dedicated USB DAC network player.

## Goals

- Raspberry Pi 4 firmware built with Buildroot and a `br2-external` tree.
- ALSA direct playback to USB DAC.
- AirPlay via `shairport-sync`.
- Roon via official RoonBridge binary package.
- BusyBox init, no `systemd` in V1.
- No PulseAudio, PipeWire, desktop, Bluetooth, or HDMI audio path.
- Tiny web controller on port `8080` for AirPlay/RoonBridge start and stop.
- Raspberry Pi Linux `rpi-6.12.y`, pinned to commit `a923c1dcd822385aaa82a4388c6a6ad92cbf9280`.

## Layout

```text
buildroot/external/
  board/rpi4-usbdac/        Raspberry Pi boot files, genimage config, rootfs overlay
  configs/                  Buildroot defconfig
  package/audioctl/         Web controller package and Go source
  package/roonbridge/       RoonBridge binary package wrapper
docs/                       Design notes and validation policy
.github/workflows/          CI image build workflow
```

## Local Build

Buildroot is tracked as a git submodule at `buildroot-src`, pinned to a specific commit from the `2025.02.x` branch.

```bash
git submodule update --init --recursive
make -C buildroot-src BR2_EXTERNAL=$PWD/buildroot/external rpi4_usbdac_defconfig
make -C buildroot-src BR2_EXTERNAL=$PWD/buildroot/external
```

The image should be generated under:

```text
buildroot-src/output/images/sdcard.img
```

To update Buildroot later:

```bash
git submodule update --remote buildroot-src
make -C buildroot-src BR2_EXTERNAL=$PWD/buildroot/external rpi4_usbdac_defconfig
```

## RoonBridge

RoonBridge is proprietary and installed from the official ARMv8 tarball:

```text
https://download.roonlabs.net/builds/RoonBridge_linuxarmv8.tar.bz2
```

Pinned SHA256:

```text
51bc2f4d0f2d79dfcf694a44e9d78469a48b68f978e474ac94e5b35d096858a9
```

The checksum is recorded in `buildroot/external/package/roonbridge/roonbridge.hash`. If Roon updates the tarball in place, Buildroot will stop with a checksum error; update the hash only after manually verifying the new official download.

## Web Controller

The `audioctl` service listens on port `8080` and exposes:

```text
GET  /api/status
POST /api/airplay/start
POST /api/airplay/stop
POST /api/roon/start
POST /api/roon/stop
```

Starting AirPlay stops RoonBridge first. Starting RoonBridge stops AirPlay first.

## Verification

Fast local check for the controller:

```bash
cd buildroot/external/package/audioctl/src
go test ./...
```

Full firmware verification requires a Raspberry Pi 4, a USB DAC, wired network, and the target Wi-Fi network.
