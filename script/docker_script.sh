#!/bin/bash

docker build -f ./Dockerfile -t StoryShift:0.0.1 .
docker run --name storyshift -e STORY_SHIFT_CONFIG_FILE=./config/test_config.yaml docker.io/library/story-shift:0.0.1