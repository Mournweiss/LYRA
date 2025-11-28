import os
from errors import ConfigError

class Config:
    TELEGRAM_TOKEN: str
    API_GATEWAY_HOST: str
    API_GATEWAY_PORT: str
    API_GATEWAY_ADDRESS: str
    MINIO_HOST: str
    MINIO_PORT: str
    MINIO_CONSOLE_PORT: str
    MINIO_ACCESS_KEY: str
    MINIO_SECRET_KEY: str
    MINIO_BUCKET: str
    MINIO_REGION: str
    MINIO_ENDPOINT: str

    @classmethod
    def load(cls) -> 'Config':
        token = os.getenv("TELEGRAM_BOT_TOKEN", "your-telegram-bot-token")
        gateway_host = os.getenv("API_GATEWAY_HOST", "api-gateway")
        gateway_port = os.getenv("API_GATEWAY_PORT", "50051")
        minio_host = os.getenv("MINIO_HOST", "minio")
        minio_port = os.getenv("MINIO_PORT", "9000")
        minio_console_port = os.getenv("MINIO_CONSOLE_PORT", "9001")
        minio_access_key = os.getenv("MINIO_ACCESS_KEY", "minioadmin")
        minio_secret_key = os.getenv("MINIO_SECRET_KEY", "minioadmin123")
        minio_bucket = os.getenv("MINIO_BUCKET", "lyra-media")
        minio_region = os.getenv("MINIO_REGION", "us-east-1")

        if not token or token == "your-telegram-bot-token":
            raise ConfigError("TELEGRAM_BOT_TOKEN environment variable is required or must be set in .env")

        if not gateway_host:
            raise ConfigError("API_GATEWAY_HOST environment variable is required or must be set in .env")

        if not gateway_port:
            raise ConfigError("API_GATEWAY_PORT environment variable is required or must be set in .env")

        if not minio_host:
            raise ConfigError("MINIO_HOST environment variable is required or must be set in .env")

        if not minio_port:
            raise ConfigError("MINIO_PORT environment variable is required or must be set in .env")

        if not minio_access_key:
            raise ConfigError("MINIO_ACCESS_KEY environment variable is required or must be set in .env")

        if not minio_secret_key:
            raise ConfigError("MINIO_SECRET_KEY environment variable is required or must be set in .env")

        if not minio_bucket:
            raise ConfigError("MINIO_BUCKET environment variable is required or must be set in .env")

        cfg = cls()
        cfg.TELEGRAM_TOKEN = token
        cfg.API_GATEWAY_HOST = gateway_host
        cfg.API_GATEWAY_PORT = gateway_port
        cfg.API_GATEWAY_ADDRESS = f"{gateway_host}:{gateway_port}"
        cfg.MINIO_HOST = minio_host
        cfg.MINIO_PORT = minio_port
        cfg.MINIO_CONSOLE_PORT = minio_console_port
        cfg.MINIO_ACCESS_KEY = minio_access_key
        cfg.MINIO_SECRET_KEY = minio_secret_key
        cfg.MINIO_BUCKET = minio_bucket
        cfg.MINIO_REGION = minio_region
        cfg.MINIO_ENDPOINT = f"{minio_host}:{minio_port}"
        return cfg
