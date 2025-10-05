from telegram.ext import Application, CommandHandler, MessageHandler, filters

def init_bot(config, handle_message, start_command):
    app = Application.builder().token(config.TELEGRAM_TOKEN).build()
    app.bot_data["config"] = config
    app.add_handler(CommandHandler("start", start_command))
    app.add_handler(MessageHandler(filters.ALL, handle_message))
    return app
