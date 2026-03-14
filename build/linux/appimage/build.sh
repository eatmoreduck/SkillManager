#!/usr/bin/env bash
# Copyright (c) 2018-Present Lea Anthony
# SPDX-License-Identifier: MIT

# Fail script on any error
set -euxo pipefail

# Define variables
APP_DIR="${APP_NAME}.AppDir"
LINUXDEPLOY_VERSION="1-alpha-20251107-1"

# Create AppDir structure
mkdir -p "${APP_DIR}/usr/bin"
mkdir -p "${APP_DIR}/usr/share/applications"
mkdir -p "${APP_DIR}/usr/share/icons/hicolor/1024x1024/apps"
cp -r "${APP_BINARY}" "${APP_DIR}/usr/bin/"
cp "${ICON_PATH}" "${APP_DIR}/"
cp "${DESKTOP_FILE}" "${APP_DIR}/"
cp "${DESKTOP_FILE}" "${APP_DIR}/usr/share/applications/$(basename "${DESKTOP_FILE}")"
cp "${ICON_PATH}" "${APP_DIR}/usr/share/icons/hicolor/1024x1024/apps/$(basename "${ICON_PATH}")"

if [[ $(uname -m) == *x86_64* ]]; then
    # Download linuxdeploy and make it executable
    wget -q -4 -N "https://github.com/linuxdeploy/linuxdeploy/releases/download/${LINUXDEPLOY_VERSION}/linuxdeploy-x86_64.AppImage"
    chmod +x linuxdeploy-x86_64.AppImage

    # Run linuxdeploy to bundle the application
    ./linuxdeploy-x86_64.AppImage --appdir "${APP_DIR}" --output appimage
else
    # Download linuxdeploy and make it executable (arm64)
    wget -q -4 -N "https://github.com/linuxdeploy/linuxdeploy/releases/download/${LINUXDEPLOY_VERSION}/linuxdeploy-aarch64.AppImage"
    chmod +x linuxdeploy-aarch64.AppImage

    # Run linuxdeploy to bundle the application (arm64)
    ./linuxdeploy-aarch64.AppImage --appdir "${APP_DIR}" --output appimage
fi

# Rename the generated AppImage regardless of the basename linuxdeploy chose
APPIMAGE_FILE="$(find . -maxdepth 1 -type f -name '*.AppImage' | head -n 1)"
if [[ -z "${APPIMAGE_FILE}" ]]; then
    echo "No AppImage file was generated" >&2
    exit 1
fi
mv "${APPIMAGE_FILE}" "${APP_NAME}.AppImage"
