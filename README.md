# clarifai-telegram

### Install
```go get github.com/saromanov/clarifai-telegram```

### Configuration
For working of this bot, you need ```Telegram token```, Clarifai ID and Clarifai Secret. You can set it as environment variables TELEGRAM_TOKEN, CLARIFAI_ID, CLARIFAI_SECRET

```go
import (
   "github.com/saromanov/clarifai-telegram"
)

client := clarifaitelegram.Client{}
client.LoadFromEnv()
client.Start()
```

Or set is not client initialization
```go
import (
   "github.com/saromanov/clarifai-telegram"
)

client := clarifaitelegram.Client{
	 ClarifaiID: "ID",
	 ClarifaiSecret: "SECRET",
	 TelegramToken: "TOKEN",

}
client.Start()
```