#!/bin/bash

echo this is bash script
systemctl is-active --quiet service
service_state = (sc query mysql80)
echo $service_state
