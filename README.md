# Telegram Presentation Access System

This application allows you to securely share Canva presentations with members of specific Telegram channels or groups.

## Key Features
- **Secure Access**: Only Telegram channel/group members can view presentations.
- **Canva Integration**: Hides the direct Canva URL using an iframe and overlay.
- **Grade Grouping**: Presentations are organized by Grade (e.g., "5-сынып", "6-сынып").
- **Deep Linking**: Clicking a link in Telegram logs the user in and redirects them *directly* to the presentation.
- **Admin Dashboard**: Easily copy presentation links to post in Telegram.

## Setup & Deployment

### 1. GitHub Repository
The code is hosted at: `https://github.com/MeirbekAshirbayev/telegram-pres`

### 2. Render.com Deployment
The application is deployed on Render.com.

**Environment Variables:**
| Key | Value |
| :--- | :--- |
| `BOT_TOKEN` | (Your Bot Token) |
| `BOT_USERNAME` | `aqquiryq_bot` |
| `AUTO_SEED` | `true` |
| `PORT` | `8080` |
| `BASE_URL` | `https://telegram-pres.onrender.com` |

### 3. Telegram Bot Setup
Ensure you have set the domain in @BotFather:
1. `/setdomain`
2. Select `@aqquiryq_bot`
3. Enter: `https://telegram-pres.onrender.com`

## Local Development (Optional)

1.  **Configuration**:
    -   Update `config.json` with your local ngrok URL.
    
2.  **Running**:
    ```bash
    go run .
    ```

## Finding Your Group ID

We have included a tool to help you find the correct ID for your group or channel.

1.  Run the tool:
    ```bash
    go run tools/get_id/main.go
    ```
2.  Send a message in your group/channel.
3.  The ID starts with `-100`. Use this in `seed.go`.
