# helsinki-guide
This repository contains the code for a Telegram bot designed to provide information about notable buildings in Helsinki. 
The bot is working at [https://t.me/HelsinkiGuide_bot](https://t.me/HelsinkiGuide_bot).

## Motivation
Helsinki boasts many fascinating buildings, but finding information about them 
can be challenging. 
This project's primary goal is to make such information more easily accessible. 
The project relies on [a dataset provided by the Helsinki City Museum](https://hri.fi/data/en_GB/dataset/helsinkilaisten-rakennusten-historiatietoja).

## Project Goals
1. Translate the dataset into English and Russian. *(pending)* ⏳
2. Create a Telegram bot to deliver this dataset to bot users. ✅
3. Allow a user to configurate his preferences. *(pending)* ⏳
4. Allow a user to search buildings per location. ✅
5. ...

## Getting Started
- Get your bot API key from [@BotFather](https://t.me/BotFather) - `BOT_TOKEN`.
- Create a new PostgreSQL database, install an [`earthdistance` extension](https://www.postgresql.org/docs/15/earthdistance.html), get a `DATABASE_URL`.
- Populate the database with data (see ["Prepare data"](#prepare-data) for details).
- Install [Docker](https://docs.docker.com/engine/).
- Build a bot container: 
```shell
make USER=<dockerhub_username> TAG=anynewtag build
```
- Apply database migrations:
```shell
DATABASE_URL=<DATABASE_URL> make migrate
```
- Run the bot:
```shell
BOT_TOKEN=<BOT_TOKEN> DATABASE_URL=<DATABASE_URL> make run
```

## Development
### Prerequisites
- Go v.1.21 or higher should be already installed.
- A bot API token `BOT_TOKEN` provided by [@BotFather](https://t.me/BotFather).
- [Docker](https://docs.docker.com/engine/) should be already installed.
- An empty Postgresql database with an installed [`earthdistance` extension](https://www.postgresql.org/docs/15/earthdistance.html).
- An environment variable `DATABASE_URL` to connect to the PostgreSQL database.
- A subscription to [the Google Translate API](https://rapidapi.com/googlecloud/api/google-translate1/) 
is required to automatically translate the source dataset into other languages.

### Installation
Open a project root directory in a console and install project dependencies:
```shell
go mod tidy
```

Apply database migrations:
```shell
DATABASE_URL=<DATABASE_URL> make migrate
```

### Start

Run the bot:
```shell
DATABASE_URL=<DATABASE_URL> go run main.go bot --token <BOT_TOKEN>
```

Get more information about available commands and options:
```shell
go run main.go --help
```

### Prepare Data

#### Translate [the source dataset](https://hri.fi/data/en_GB/dataset/helsinkilaisten-rakennusten-historiatietoja)

This command will create a new file `translated.xlsx` where a `Lauttasaari`
sheet will be partially translated into English.
```shell
go run main.go translate --api-key <your Google Translate API key> --sheet Lauttasaari input_dataset.xlsx translated.xlsx
```

#### Populate the database

Transfer the data from `xlsx` files to the database:
```shell
go run main.go populate --dburl ${DatabaseURL} --sheet Lauttasaari fi.xlsx en.xlsx ru.xlsx
```

### Tests

Run the project tests: 
```shell
make test
```

## Acknowledgements
Source: History of buildings in Helsinki. The maintainer of the dataset is Helsingin kulttuurin ja vapaa-ajan toimiala / Kaupunginmuseo and the original author is Tmi Hilla Tarjanne. The dataset has been downloaded from Helsinki Region Infoshare service on 2023-10-22 18:00:08.977295 under the license Creative Commons Attribution 4.0. 

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
