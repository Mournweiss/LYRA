import logging
from config import Config
from errors import ConfigError
from handlers.message import handle_message
from bot.commands import start_command
from bot.init import init_bot

def main():
    logging.basicConfig(level=logging.INFO)

    try:
        config = Config.load()

    except ConfigError as e:
        logging.error(f"Configuration error: {e}")
        exit(1)
        
    app = init_bot(config, handle_message, start_command)
    app.run_polling()

if __name__ == "__main__":
    main()
