#pragma once
#include <stdexcept>
#include <string>

class BaseError : public std::runtime_error {
public:
    explicit BaseError(const std::string& msg) : std::runtime_error(msg) {}
};

class ConfigError : public BaseError {
public:
    explicit ConfigError(const std::string& msg) : BaseError(msg) {}
};

class UpstreamError : public BaseError {
public:
    explicit UpstreamError(const std::string& msg) : BaseError(msg) {}
};

class StartupError : public BaseError {
public:
    explicit StartupError(const std::string& msg) : BaseError(msg) {}
};
