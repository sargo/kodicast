#!/bin/bash
sleep 20
/home/osmc/bin/kodicast -log-kodi -log-player -loglevel info 2>&1 | tee -a /var/log/kodicast.log

