# EstreiaBot

A simple Telegram bot written in Go that alert the user when a Tv Show from his choice releases a new Season. With this bot, users can add Tv Shows to their watching list and receive notifications when a new season is released. This project uses the TMDB API.

## Features

- **Add Tv Show to Your List**: Users can add their favorite Tv Shows series to a personal list.
- **Season Notifications**: Receive alerts when a new season is released for Tv Shows in your list.

## How to test

If you have a telegram account, you can use this bot sending a message to https://t.me/EstreiaBot

## Prerequisites

Before running the bot, ensure you have the following installed:

- [Go](https://go.dev/dl/)
- [Telegram Bot Token](https://core.telegram.org/bots#botfather): You'll need to create a bot on Telegram and get your unique API token from BotFather.
- [TMDB API Key](https://developer.themoviedb.org/reference/intro/getting-started): You'll need to get a TMDB API Key

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/joaovs2004/EstreiaBot
cd EstreiaBot/
```

### 2. Set the Telegram API Token on .env like this
```
TELEGRAM_TOKEN="Your Token here"
TMDB_API_KEY="Your API key here"
```

### 3. Run the project

```bash
go run .
```