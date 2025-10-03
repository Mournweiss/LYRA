import os

class ConfigError(Exception):
    pass

class Config:
    TELEGRAM_TOKEN: str
    API_GATEWAY_ADDRESS: str

    @classmethod
    def load(cls) -> 'Config':
        token = os.getenv("TELEGRAM_BOT_TOKEN")

        if not token:
            raise ConfigError("TELEGRAM_BOT_TOKEN environment variable is required.")

        gateway_addr = os.getenv("API_GATEWAY_ADDRESS")

        if not gateway_addr:
            raise ConfigError("API_GATEWAY_ADDRESS environment variable is required.")

        cfg = cls()
        cfg.TELEGRAM_TOKEN = token
        cfg.API_GATEWAY_ADDRESS = gateway_addr
        return cfg
