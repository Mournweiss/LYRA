#include "config.h"
#include "errors.h"
#include <cstdlib>

Config Config::Load() {
    Config cfg;
    const char* port = std::getenv("WHISPER_SERVICE_PORT");
    std::string default_port = "50052";
    
    if (!port || std::string(port).empty()) {
        cfg.service_port = default_port;
    } else {
        cfg.service_port = port;
    }

    if (cfg.service_port.empty()) {
        throw ConfigError("WHISPER_SERVICE_PORT environment variable is required or must be set in .env");
    }

    return cfg;
}
