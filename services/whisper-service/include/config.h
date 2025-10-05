#pragma once
#include <string>
#include "errors.h"

class Config {
public:
    std::string service_port;
    std::string service_host;
    std::string minio_host;
    std::string minio_port;
    std::string minio_access_key;
    std::string minio_secret_key;
    std::string minio_bucket;
    std::string minio_endpoint;
    static Config Load();
};
