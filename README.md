# helsinki-guide
This repository contains the code for a Telegram bot designed to provide information about notable buildings in Helsinki. 
The bot is working at [https://t.me/HelsinkiGuide_bot](https://t.me/HelsinkiGuide_bot).

Please note that this project is still a work in progress.

## Motivation
Helsinki boasts many fascinating buildings, but finding information about them 
can be challenging. 
This project's primary goal is to make such information more easily accessible. 
The project relies on [a dataset provided by the Helsinki City Museum](https://hri.fi/data/en_GB/dataset/helsinkilaisten-rakennusten-historiatietoja).

## Project Goals
1. Translate the dataset into English and Russian. *(pending)* ⏳
2. Create a Telegram bot to deliver this dataset to bot users. ✅
3. Allow a user to configurate his preferences. (TO DO) 🏡
4. Allow a user to search per location. (TO DO) 🏡
5. ...

## Getting Started
- Get your bot API key from [@BotFather](https://t.me/BotFather) - `BotAPIKey`.
- Create a new PostgreSQL database, get a `DatabaseURL`.
- Populate the database with data (see "Prepare data" for details).
- Install (Docker)[https://docs.docker.com/engine/].
- Build a bot container: 
```bash
$ make build
```
- Set `BotAPIKey` and `DatabaseURL` as environment variables.
- Run the bot:
```bash
$ make run
```

## Development
### Prerequisites
- Go v.1.21 or higher should be already installed.
- Docker should be already installed.
- A subscription to [the Google Translate API](https://rapidapi.com/googlecloud/api/google-translate1/) 
is required to automatically translate the source dataset into other languages.

### Installation

Open a project root directory in a console and install project dependencies:
```shell
go mod tidy
```

### Start

Translate the source dataset into English:
```shell
go run main.go translate --api-key <your Google Translate API key> --sheet Pohjois-Haaga input_dataset.xlsx translated.xlsx
```
This command will create a new file `translated.xlsx` where a `Pohjois-Haaga`
sheet will be translated into English.

Get more information about available commands and options:
```shell
go run main.go --help
```

### Tests

Run the project tests: 
```shell
go test ./...
```

## Acknowledgements
Source: History of buildings in Helsinki. The maintainer of the dataset is Helsingin kulttuurin ja vapaa-ajan toimiala / Kaupunginmuseo and the original author is Tmi Hilla Tarjanne. The dataset has been downloaded from Helsinki Region Infoshare service on 2023-10-22 18:00:08.977295 under the license Creative Commons Attribution 4.0. 

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details