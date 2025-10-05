import os
from errors import ConfigError

class Config:
    TELEGRAM_TOKEN: str
    API_GATEWAY_HOST: str
    API_GATEWAY_PORT: str
    API_GATEWAY_ADDRESS: str

    @classmethod
    def load(cls) -> 'Config':
        token = os.getenv("TELEGRAM_BOT_TOKEN", "your-telegram-bot-token")
        gateway_host = os.getenv("API_GATEWAY_HOST", "api-gateway")
        gateway_port = os.getenv("API_GATEWAY_PORT", "50051")

        if not token or token == "your-telegram-bot-token":
            raise ConfigError("TELEGRAM_BOT_TOKEN environment variable is required or must be set in .env")
            
        if not gateway_host:
            raise ConfigError("API_GATEWAY_HOST environment variable is required or must be set in .env")

        if not gateway_port:
            raise ConfigError("API_GATEWAY_PORT environment variable is required or must be set in .env")

        cfg = cls()
        cfg.TELEGRAM_TOKEN = token
        cfg.API_GATEWAY_HOST = gateway_host
        cfg.API_GATEWAY_PORT = gateway_port
        cfg.API_GATEWAY_ADDRESS = f"{gateway_host}:{gateway_port}"
        return cfg
