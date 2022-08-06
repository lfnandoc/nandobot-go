# nandobot-go

nandobot-go is a Golang interpretation of my previous project Nandobot.

Instead of being a Discord bot, it uses Discord Webhooks to send the discord messages.

Also new:
- Generates an image of the player's items stitched together and uploads it to Imgur
- Has an API for adding new players and getting the last match info
- Stores all matches on db

![image](https://user-images.githubusercontent.com/82987034/177017838-a65f8cf0-df37-48ad-85bf-b10f84be0489.png)

New Config.json:
```
{
    "riotApiKey": "RGAPI-xxxxxxxxxxxxxxxxxxxxxxxx", // your riot games api key
    "ImgurClientID": "xxx", // imgur client id
    "ImgurClientSecret": "xx", // imgur client secret
    "DiscordWebhook": "https://discordapp.com/api/webhooks/xxxxxxxxxx" // discord webhook
}
```
