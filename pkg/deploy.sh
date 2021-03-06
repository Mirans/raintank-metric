#!/bin/bash

# Find the directory we exist within
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ -z ${PACKAGECLOUD_REPO} ] ; then
  echo "The environment variable PACKAGECLOUD_REPO must be set."
  exit 1
fi

package_cloud push ${PACKAGECLOUD_REPO}/ubuntu/trusty ${DIR}/artifacts/*.deb
git push --tags
