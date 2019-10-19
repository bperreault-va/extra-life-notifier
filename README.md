# extra-life-notifier

Extra Life Notifier will ping your slack or discord channel whenever a member of your Extra Life team received a donation. You'll also get a ping whenever someone joins the team!

## Executable

Run an executable file for your system. Supports:
- windows/386
- darwin/amd64

Find your Extra Life Team ID in the URL of your [Team Page](https://imgur.com/ibu50IB)

Follow the instructions in [this article](https://slack.com/intl/en-ca/help/articles/115005265063-incoming-webhooks-for-slack) to create a Slack Incoming Webhook URL.

Follow the instructions in [this article](https://support.discordapp.com/hc/en-us/articles/228383668-Intro-to-Webhooks) to create a Discord Webhook URL.

After adding either Slack or Discord, choose "Start server". A short message will be sent to your webhook URLs to test them. If everything looks good, then new donations and participants will trigger a ping to your app!


## Build

After pulling the source files, `main.go` can be run to start a server. Just add your Extra Life Team ID and Slack and Discord Webhook URLS to the corresponding variables in `main.go`.

An executable can be built for your operating system by running `go build -o extralife exec/main.go` from root. If you're on a Windows OS, replace `extralife` with `extralife.exe`.

## Hosting on Google Cloud Platform

It's quick and easy to host this script on Google Cloud Platform. Charges may apply, but they should be minimal.
1. Create a new project at `https://console.cloud.google.com`. Log in with a Google account or create a new one.
2. Download and install Cloud SDK, then download and install Go at `https://cloud.google.com/appengine/docs/flexible/go/download`.
3. Add your Extra Life Team ID and webhook urls in `main.go`.
4. From the extra-life-notifier project root:
- Run `gcloud auth login` and follow instructions
- Run `gcloud config set project [YOUR PROJECT ID]`
- Run `gcloud app deploy` and follow instructions.
