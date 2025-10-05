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

class MinioError : public BaseError {
public:
    explicit MinioError(const std::string& msg) : BaseError(msg) {}
};

class WhisperError : public BaseError {
public:
    explicit WhisperError(const std::string& msg) : BaseError(msg) {}
};

class TaskError : public BaseError {
public:
    explicit TaskError(const std::string& msg) : BaseError(msg) {}
};
