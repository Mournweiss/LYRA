#include "config.h"
#include "errors.h"
#include <cstdlib>
#include <string>

Config Config::Load() {
    Config cfg;
    const char* port = std::getenv("WHISPER_SERVICE_PORT");
    const char* domain = std::getenv("WHISPER_SERVICE_DOMAIN");
    std::string default_port = "50052";
    std::string default_domain = "whisper-service";

    cfg.service_port = (port && std::string(port).length() > 0) ? port : default_port;
    cfg.service_domain = (domain && std::string(domain).length() > 0) ? domain : default_domain;

    if (cfg.service_port.empty()) {
        throw ConfigError("WHISPER_SERVICE_PORT environment variable is required or must be set in .env");
    }
    if (cfg.service_domain.empty()) {
        throw ConfigError("WHISPER_SERVICE_DOMAIN environment variable is required or must be set in .env");
    }

    return cfg;
}
