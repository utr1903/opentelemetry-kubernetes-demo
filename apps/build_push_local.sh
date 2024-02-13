#!/bin/bash

# Get commandline arguments
while (( "$#" )); do
  case "$1" in
    --docker-username)
      dockerUsername="${2}"
      shift
      ;;
    --platform)
      platform="$2"
      shift
      ;;
    --language)
      language="$2"
      shift
      ;;
    --project)
      project="${2}"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

# Docker platform
if [[ $platform == "" ]]; then
  # Default is amd
  platform="amd64"
else
  if [[ $platform != "amd64" && $platform != "arm64" ]]; then
    echo "Platform [--platform] can either be 'amd64' or 'arm64'."
    exit 1
  fi
fi

# Programming language
if [[ $language != "golang" ]]; then
  echo "Currently supported languages [--language] are 'golang'."
  exit 1
fi

httpserverImageName="${project}-httpserver-${language}:latest"
kafkaconsumerImageName="${project}-kafkaconsumer-${language}:latest"
simulatorImageName="${project}-simulator-${language}:latest"
latencymanagerImageName="${project}-latencymanager-${language}:latest"

####################
### Build & Push ###
####################

# httpserver
docker build \
  --platform "linux/${platform}" \
  --tag "${dockerUsername}/${httpserverImageName}" \
  --build-arg="APP_NAME=httpserver" \
  "./${language}/."
docker push "${dockerUsername}/${httpserverImageName}"

# kafkaconsumer
docker build \
  --platform "linux/${platform}" \
  --tag "${dockerUsername}/${kafkaconsumerImageName}" \
  --build-arg="APP_NAME=kafkaconsumer" \
  "./${language}/."
docker push "${dockerUsername}/${kafkaconsumerImageName}"

# simulator
docker build \
  --platform "linux/${platform}" \
  --tag "${dockerUsername}/${simulatorImageName}" \
  --build-arg="APP_NAME=simulator" \
  "./${language}/."
docker push "${dockerUsername}/${simulatorImageName}"

# latencymanager
docker build \
  --platform "linux/${platform}" \
  --tag "${dockerUsername}/${latencymanagerImageName}" \
  --build-arg="APP_NAME=latencymanager" \
  "./${language}/."
docker push "${dockerUsername}/${latencymanagerImageName}"
