# SEU (Single Event Upset)

Watch this and you will understand the purpose of this code


[The Universe is Hostile to Computers](https://www.youtube.com/watch?v=AaZ_RSt0KP8)


Code create a big (10MB) array and every minute check him for a changes. 
If changes detected it logged.

# Usage
Logging is carried out using a telegram bot, or just in output.txt if failed send to bot.
Run code with flags:

 - chat_id - it can be your telegram id or chat id
 -  token - telegram bot token ([how to create bot](https://core.telegram.org/bots/features#creating-a-new-bot))
    
    ./seu --chat_id $CHAT_ID --token $TOKEN
