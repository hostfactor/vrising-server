#!/bin/bash

$HOME/renfield &

Xvfb :0 -screen 0 1024x768x16 &

DISPLAY=:0.0 wine64 $HOME/vrising/VRisingServer.exe -persistentDataPath "$SAVE_DIR" -logFile "$LOG_FILE" -serverName "$SERVER_NAME" -saveName "$SAVE_NAME" "$@"
