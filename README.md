# Notify.is Google Cloud Run Deployment

![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/oliverproud/notify.is) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/oliverproud/notify.is-go)

Checks username availability, sends success email via SendGrid, updates database

[Notify.is website repository](https://github.com/oliverproud/notify.is)

**Todo**:
1. Fix username validation as different services require usernames to be in different formats.
2. Fix issue with GitHub check giving a false positive when the username that it checked with was of the wrong format.
3. Fix issue with Instagram check sometimes resulting in a false negative.
