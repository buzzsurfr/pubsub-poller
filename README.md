# pubsub-poller
Google Cloud Pub/Sub poller

This was made for one of my classes at ECU.

## Installing

See [releases](https://github.com/buzzsurfr/pubsub-poller/releases) for different builds. There's also a docker image at `ghcr.io/buzzsurfr/pubsub-poller:v0.0.0`.

## Usage

```
Polls gcloud pubsub subscriptions and outputs messages to stdout

Usage:
  pubsub-poller [FLAGS] SUBSCRIPTION... [flags]

Flags:
  -h, --help                help for pubsub-poller
  -p, --project-id string   Project ID
```

### Running in Docker container

First, make sure you have an application login by running (and login):

```
gcloud auth application-default login
```

Then, run the docker container:

```
docker run -it --rm -v $HOME/.config/gcloud:/root/.config/gcloud ghcr.io/buzzsurfr/pubsub-poller:v0.0.0
```

### Sample Output

The output containes the time (since launch), the actual message (mine is just the resource ID), the pubsub message ID, the project ID, and the subscription name.

```
INFO[0002] 5671636300726272                              message-id=2161735373414595 project-id=seng-6285-spring21 subscription=subject-deleted-sub
INFO[0002] 5671636300726272                              message-id=2161751563016368 project-id=seng-6285-spring21 subscription=subject-deleted-sub
INFO[0002] 5671636300726272                              message-id=2161752335379751 project-id=seng-6285-spring21 subscription=subject-deleted-sub
INFO[0002] 5671636300726272                              message-id=2161642375748982 project-id=seng-6285-spring21 subscription=subject-deleted-sub
```

### Restful Bots Homework

Fellow students, use this to check for all the topics in Homework 4. Make sure your Google Cloud Project ID is set as `GCP_PROJECT`

```
docker run -it --rm -v $HOME/.config/gcloud:/root/.config/gcloud ghcr.io/buzzsurfr/pubsub-poller:v0.0.0 -p $GCP_PROJECT course-created-sub course-updated-sub course-deleted-sub subject-created-sub subject-updated-sub subject-deleted-sub textbook-created-sub textbook-updated-sub textbook-deleted-sub
```