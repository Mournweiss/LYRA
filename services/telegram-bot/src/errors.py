class BotError(Exception):
    pass

class ConfigError(BotError):
    pass

class GRPCError(BotError):
    pass

class MinioError(BotError):
    pass

class GatewayError(BotError):
    pass

class ValidationError(BotError):
    pass
