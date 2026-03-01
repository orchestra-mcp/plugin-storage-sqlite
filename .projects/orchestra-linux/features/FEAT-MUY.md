---
created_at: "2026-02-28T02:54:40Z"
description: 'Implement services.voice TTS tools on Linux. Primary: Speech Dispatcher (libspeechd) — meta-engine that wraps espeak-ng, festival, piper. spd_open(), spd_say() with SPD_TEXT priority. Secondary: Piper neural TTS (high quality offline) — invoke piper CLI with --model path, pipe WAV output to PipeWire/PulseAudio via GStreamer (gst-launch-1.0 filesrc → wavparse → autoaudiosink). Voice list from spd_list_synthesis_voices(). Provider TTS (ElevenLabs, OpenAI TTS, Google Cloud TTS) via HTTP API calls using libsoup3. Calls tts_speak, tts_speak_provider, tts_list_voices, tts_stop tools.'
id: FEAT-MUY
priority: P2
project_id: orchestra-linux
status: backlog
title: Voice TTS via Speech Dispatcher + Piper
updated_at: "2026-02-28T02:54:40Z"
version: 0
---

# Voice TTS via Speech Dispatcher + Piper

Implement services.voice TTS tools on Linux. Primary: Speech Dispatcher (libspeechd) — meta-engine that wraps espeak-ng, festival, piper. spd_open(), spd_say() with SPD_TEXT priority. Secondary: Piper neural TTS (high quality offline) — invoke piper CLI with --model path, pipe WAV output to PipeWire/PulseAudio via GStreamer (gst-launch-1.0 filesrc → wavparse → autoaudiosink). Voice list from spd_list_synthesis_voices(). Provider TTS (ElevenLabs, OpenAI TTS, Google Cloud TTS) via HTTP API calls using libsoup3. Calls tts_speak, tts_speak_provider, tts_list_voices, tts_stop tools.
