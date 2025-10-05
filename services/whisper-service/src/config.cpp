#include "config.h"
#include "errors.h"
#include <cstdlib>
#include <string>

Config Config::Load() {
    Config cfg;
    const char* port = std::getenv("WHISPER_SERVICE_PORT");
    const char* host = std::getenv("WHISPER_SERVICE_HOST");
    std::string default_port = "50052";
    std::string default_host = "whisper-service";

    cfg.service_port = (port && std::string(port).length() > 0) ? port : default_port;
    cfg.service_host = (host && std::string(host).length() > 0) ? host : default_host;

    if (cfg.service_port.empty()) {
        throw ConfigError("WHISPER_SERVICE_PORT environment variable is required or must be set in .env");
    }
    if (cfg.service_host.empty()) {
        throw ConfigError("WHISPER_SERVICE_HOST environment variable is required or must be set in .env");
    }

    return cfg;
}
