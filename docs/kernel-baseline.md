# Kernel And Runtime Baseline

This firmware targets Raspberry Pi 4 as a dedicated USB DAC network player.

The V1 kernel baseline is Raspberry Pi Linux `rpi-6.12.y`, pinned to commit `a923c1dcd822385aaa82a4388c6a6ad92cbf9280`, with Buildroot Linux headers set to `6.12`.

The first version intentionally uses a regular low-latency/preemptible kernel instead of PREEMPT_RT. The main goal is a quiet and predictable playback appliance: fewer drivers, fewer services, stable ALSA routing, and less background I/O.

Keep enabled:

- Raspberry Pi 4 boot, storage, USB host, Ethernet, and Wi-Fi support.
- ALSA core and USB Audio Class support through `snd-usb-audio`.
- High resolution timers and kernel preemption.
- CPU frequency control, with runtime policy set to `performance` during playback.
- `ext4` and `vfat` for root and boot partitions.

Avoid in V1:

- `PREEMPT_RT`.
- `systemd`.
- PulseAudio, PipeWire, desktop shells, Bluetooth, HDMI audio, and broad debug/tracing options.
- Unrelated sound card, camera, printer, and display features.

The practical tuning order is: stable USB DAC detection, ALSA default device, service mutual exclusion, network reliability, then kernel experiments.
