# Notify.is Google Cloud Run Deployment

![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/oliverproud/notify.is) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/oliverproud/notify.is-go)

Checks username availability, sends success email via SendGrid, updates database

[Notify.is website repository](https://github.com/oliverproud/notify.is)

**Todo**:
1. ~~Fix username validation as different services require usernames to be in different formats.~~ [Fixed - [commit](https://github.com/oliverproud/notify.is/commit/fe95bb4a45a47aa5b72bd918eef83490954691cc)]
2. ~~Fix issue with GitHub check giving a false positive when the username that it checked with was of the wrong format.~~ [Fixed - same as above]
3. ~~Fix issue with Instagram check sometimes resulting in a false negative.~~ [Fixed - [commit](https://github.com/oliverproud/notify.is-gcloud/commit/1fc17d8ba91373e5334f46566fafb1a87484d89b)]
