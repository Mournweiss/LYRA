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

    const char* minio_host = std::getenv("MINIO_HOST");
    const char* minio_port = std::getenv("MINIO_PORT");
    const char* minio_access_key = std::getenv("MINIO_ACCESS_KEY");
    const char* minio_secret_key = std::getenv("MINIO_SECRET_KEY");
    const char* minio_bucket = std::getenv("MINIO_BUCKET");
    cfg.minio_host = minio_host ? minio_host : "minio";
    cfg.minio_port = minio_port ? minio_port : "9000";
    cfg.minio_access_key = minio_access_key ? minio_access_key : "minioadmin";
    cfg.minio_secret_key = minio_secret_key ? minio_secret_key : "minioadmin123";
    cfg.minio_bucket = minio_bucket ? minio_bucket : "lyra-media";
    cfg.minio_endpoint = cfg.minio_host + std::string(":") + cfg.minio_port;

    if (cfg.service_port.empty()) {
        throw ConfigError("WHISPER_SERVICE_PORT environment variable is required or must be set in .env");
    }
    if (cfg.service_host.empty()) {
        throw ConfigError("WHISPER_SERVICE_HOST environment variable is required or must be set in .env");
    }
    if (cfg.minio_host.empty()) {
        throw ConfigError("MINIO_HOST environment variable is required or must be set in .env");
    }
    if (cfg.minio_port.empty()) {
        throw ConfigError("MINIO_PORT environment variable is required or must be set in .env");
    }
    if (cfg.minio_access_key.empty()) {
        throw ConfigError("MINIO_ACCESS_KEY environment variable is required or must be set in .env");
    }
    if (cfg.minio_secret_key.empty()) {
        throw ConfigError("MINIO_SECRET_KEY environment variable is required or must be set in .env");
    }
    if (cfg.minio_bucket.empty()) {
        throw ConfigError("MINIO_BUCKET environment variable is required or must be set in .env");
    }

    return cfg;
}
