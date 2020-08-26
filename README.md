# go-slack-bot
Slack bot written in golang. Bot can easily understand and respond to many greetings, and can answer some general knowledge questions like - "Who is the president of USA?"
# Setup
1. Get a token for slack bot, wit.ai and Wolfram Alpha API OR contact me for some access tokens, and to be invited to a slack channel
2. Replace your tokens with the placeholders in start_bot.sh.
3. Run ```./start_bot.sh```. If you are on windows use git bash or manually set environment variables using ```SLACK_TOKEN=YOUR_TOKEN WIT_TOKEN=YOUR_TOKEN WOLFRAM_TOKEN=YOUR_TOKEN go run bot.go```

