class BotError(Exception):
    pass

class ConfigError(BotError):
    pass

class GRPCError(BotError):
    pass
