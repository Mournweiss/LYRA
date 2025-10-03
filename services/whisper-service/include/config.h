#pragma once
#include <string>
#include <stdexcept>

class ConfigError : public std::runtime_error {
public:
    explicit ConfigError(const std::string& msg) : std::runtime_error(msg) {}
};

class Config {
public:
    std::string service_port;
    
    static Config Load();
};
