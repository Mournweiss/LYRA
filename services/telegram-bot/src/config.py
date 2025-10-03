import os
from errors import ConfigError

class Config:
    TELEGRAM_TOKEN: str
    API_GATEWAY_ADDRESS: str

    @classmethod
    def load(cls) -> 'Config':
        token = os.getenv("TELEGRAM_BOT_TOKEN", "your-telegram-bot-token")
        gateway_addr = os.getenv("API_GATEWAY_ADDRESS", "api-gateway:50051")

        if not token or token == "your-telegram-bot-token":
            raise ConfigError("TELEGRAM_BOT_TOKEN environment variable is required or must be set in .env")
            
        if not gateway_addr or gateway_addr == "api-gateway:${API_GATEWAY_PORT}":
            raise ConfigError("API_GATEWAY_ADDRESS environment variable is required or must be set in .env")

        cfg = cls()
        cfg.TELEGRAM_TOKEN = token
        cfg.API_GATEWAY_ADDRESS = gateway_addr
        return cfg
