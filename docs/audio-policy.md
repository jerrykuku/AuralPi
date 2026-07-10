# Audio Policy

The system exposes two network playback modes:

- AirPlay through `shairport-sync`.
- Roon through official `RoonBridge`.

Both services target ALSA directly. The firmware does not include PulseAudio, PipeWire, software mixing, or a desktop audio server.

V1 uses a strict service-level policy: starting AirPlay stops RoonBridge first, and starting RoonBridge stops AirPlay first. This avoids two playback stacks competing for the same USB DAC.

The default ALSA device is pinned in `/etc/asound.conf` to card `1`, which is the common Raspberry Pi layout when HDMI audio is disabled and a USB DAC is present. A later version should replace this with a boot-time USB DAC discovery script that generates the ALSA default dynamically.

Validation checklist:

- `aplay -l` shows the USB DAC.
- `speaker-test -D default` plays through the USB DAC.
- AirPlay playback stops RoonBridge before opening ALSA.
- Starting RoonBridge from the web UI stops AirPlay before launching RoonBridge.
- Long playback sessions do not produce underruns or service restarts.
