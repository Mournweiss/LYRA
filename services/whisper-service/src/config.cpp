#include "config.h"
#include <cstdlib>

Config Config::Load() {
    Config cfg;
    const char* port = std::getenv("WHISPER_SERVICE_PORT");
    
    if (!port || std::string(port).empty()) {
        throw ConfigError("WHISPER_SERVICE_PORT environment variable is required");
    }

    cfg.service_port = port;
    return cfg;
}
