# Telegram Bot API Documentation Assistant

This Telegram bot provides quick access to the Telegram Bot API documentation. Users can search for API methods and types using inline queries, and the bot will return detailed information about the specified API elements.

### Demo Bot: [@BotApiDocsBot](https://t.me/BotApiDocsBot)

## Features
- Inline search for API methods and types.
- Detailed descriptions, return types, and required fields for each method and type.
- Easy access to the official documentation.

## Prerequisites
- Go version 1.23 or higher.
- A Telegram bot token obtainable from [BotFather](https://core.telegram.org/bots#botfather).

## Installation

### 1. Install Go
Follow the instructions to install Go on your system: [Go Installation Guide](https://golang.org/doc/install).

<details>
<summary>Easy Way:</summary>

```shell
git clone https://github.com/udhos/update-golang dlgo && cd dlgo && sudo ./update-golang.sh && source /etc/profile.d/golang_path.sh
```

Exit the terminal and open it again to check the installation.
</details>

Verify the installation by running:

```shell
go version
```

### 2. Clone the repository

```shell
git clone https://github.com/AshokShau/BotApiDocs&& cd BotApiDocs
```

### 3. Set up the environment

Copy the sample environment file and edit it as needed:

```shell
cp sample.env .env
vi .env
```

### 4. Build the project

```shell
go build
```

### 5. Run the project

```shell
./BotApiDocs
```

## Usage

1. **Start a chat** with your bot on Telegram. Once the bot is running, you can search for API methods and types.
2. Use the inline query feature by typing `@YourBotUsername <your_query>` to search for methods or types.
3. The bot will return relevant results with detailed descriptions.

## Contributing
<details>
<summary>Contribution Guidelines</summary>

Contributions are welcome! Here's how you can help:

1. **Fork the repository**.
2. **Clone your forked repository** to your local machine.
    ```shell
    git clone https://github.com/your-username/BotApiDocs.git
    cd BotApiDocs
    ```
3. **Create a new branch** for your changes.
    ```shell
    git checkout -b feature-branch
    ```
4. **Make your changes** and commit them with a descriptive message.
    ```shell
    git add .
    git commit -m "Description of your changes"
    ```
5. **Push to your branch** and submit a pull request.

Please ensure your code follows the project's coding standards and includes appropriate tests.
</details>
        
## License

This project is licensed under the MIT License—see the [LICENSE](/LICENSE) file for details.

## Contact

[![Telegram](https://img.shields.io/badge/Telegram-Channel-blue.svg)](https://t.me/FallenProjects)  
[![Telegram](https://img.shields.io/badge/Telegram-Chat-blue.svg)](https://t.me/AshokShau)


## Acknowledgments
- **[Ashok Shau](https://github.com/AshokShau)**: For creating and maintaining this [project](https://github.com/AshokShau/BotApiDocs), which provides a solid foundation for building Telegram bots.

- **[PaulSonOfLars](https://github.com/PaulSonOfLars)**: For the invaluable [GoTgBot](https://github.com/PaulSonOfLars/gotgbot) library, which simplifies Telegram bot development in Go, and for the [API specification](https://github.com/PaulSonOfLars/telegram-bot-api-spec/raw/main/api.json) that serves as a reference for bot methods and types.
