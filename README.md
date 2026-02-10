# Telegram Presentation Access System

This application allows you to securely share Canva presentations with members of specific Telegram channels or groups.

## Key Features
- **Secure Access**: Only Telegram channel/group members can view presentations.
- **Canva Integration**: Hides the direct Canva URL using an iframe and overlay.
- **Grade Grouping**: Presentations are organized by Grade (e.g., "5-сынып", "6-сынып").
- **Deep Linking**: Clicking a link in Telegram logs the user in and redirects them *directly* to the presentation.
- **Admin Dashboard**: Easily copy presentation links to post in Telegram.

## Prerequisites

1.  **Go**: Ensure Go is installed on your system.
2.  **Telegram Bot**:
    -   Create a new bot via [@BotFather](https://t.me/BotFather).
    -   Get the **Bot Token**.
    -   Set the bot domain using `/setdomain` in BotFather to your public URL (or ngrok URL for testing).
3.  **Telegram Group/Channel**:
    -   Create a **Supergroup** (for Topics support) or a Channel.
    -   Add your bot as an **Administrator**.
    -   Get the **Group/Channel ID**.

## Finding Your Group ID

We have included a tool to help you find the correct ID for your group or channel.

1.  Run the tool:
    ```bash
    go run tools/get_id/main.go
    ```
2.  Send a message in your group/channel.
3.  The tool will print the ID (starts with `-100`). Use this ID in your database/seed.go.

## Setup

1.  **Configuration**:
    -   Open `config.json`.
    -   `bot_token`: Your Telegram Bot Token.
    -   `bot_username`: Your bot's username (without @).
    -   `admin_ids`: Your Telegram User ID (for future admin features).
    -   `base_url`: The public URL where this app is hosted (e.g., `https://your-domain.com` or `https://<ngrok-id>.ngrok-free.app`). *Crucial for Telegram Login.*
    -   `port`: The port to run on (default `:8080`).

2.  **Database Seeding**:
    -   To add presentations, update `seed.go` with your titles, Canva URLs, Grade names, and your **Group ID**.
    -   Run:
        ```bash
        go run . -seed
        ```
    -   This will create/update `presentations.db`.

## Running

1.  Start the server:
    ```bash
    go run .
    ```
2.  Access the admin dashboard at your `base_url`.
3.  Copy links and post them to your Telegram group topics.
4.  Users in that group can click the links to view the presentations.
