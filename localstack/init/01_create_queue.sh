#!/bin/bash
set -e

awslocal sqs create-queue --queue-name events
echo "SQS queue 'events' created"
