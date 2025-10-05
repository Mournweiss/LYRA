#pragma once
#include <string>
#include "errors.h"

class Config {
public:
    std::string service_port;
    std::string service_host;
    static Config Load();
};
