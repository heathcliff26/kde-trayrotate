#!/bin/bash

echo "${1}" >> "${KDE_TRAYROTATE_DOCTOR_LOG}"

if [ "${KDE_TRAYROTATE_DOCTOR_FAIL}" == "true" ]; then
    exit 1
fi

if [ -n "${KDE_TRAYROTATE_DOCTOR_DATA}" ]; then
    cat "${KDE_TRAYROTATE_DOCTOR_DATA}"
fi
