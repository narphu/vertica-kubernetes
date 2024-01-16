#!/bin/bash

# (c) Copyright [2021-2023] Open Text.
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Image utilities to be sourced into various bash scripts

LAST_RELEASED_IMAGE="24.1.0"
# Next two variables define the version that is built nightly from the server
# master branch. Update this as the server repo changes the version.
NIGHTLY_MAJOR=24
NIGHTLY_MINOR=2

function print_vertica_k8s_img
{
    imageName=$1
    major=$2
    minor=$3
    patch=$4
    local VERTICA_REPO="vertica"
    echo "${VERTICA_REPO}/$imageName:$major.$minor.$patch-0"
}

function get_rpm_version 
{
    local SCRIPT_DIR REPO_DIR
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
    REPO_DIR=$(dirname $SCRIPT_DIR)
    # Find the RPM version that's download and built for the CI
    grep 'VERTICA_CE_URL:' $REPO_DIR/.github/actions/download-rpm/action.yaml | cut -d':' -f3 | cut -d'/' -f5 | cut -d'-' -f2
}

function determine_image_version() {
    local TARGET_IMAGE=$1

    # Extract out the tag from the image.
    IFS=':' read image tag <<< "$TARGET_IMAGE"

    if [[ -z "$tag" ]]
    then
       # No tag found. Assume latest, so pick the last released image
       echo ${LAST_RELEASED_IMAGE}
       return
    fi

    IFS='.' read major minor patch <<< "$tag"

    # If we were able to extract only digits for major/minor, then the tag was
    # in fact a version.
    if [[ $major =~ ^[0-9]+$ && $minor =~ ^[0-9]+$ ]]
    then
        echo "$major.$minor.0"
        return
    fi

    # No able to figure out the version from the tag.  If the image repo is
    # dockerhub, then we assume we are running with the nightly build. So, we
    # return an image based on the nightly version. This must come from the private
    # repo in case the base version isn't released yet.
    if [[ $TARGET_IMAGE == docker.io/* ]]
    then
        echo "$NIGHTLY_MAJOR.$NIGHTLY_MINOR.0"
        return
    fi

    # We assume we are running with an image built in this CI that used the public
    # RPM. This is true for PRs or running off of main
    IFS='.' read major minor patch <<< "$(get_rpm_version)"
    echo "$major.$minor.$patch"
    return
}
