#pragma once
#include <string>
#include "errors.h"

class Config {
public:
    std::string service_port;
    
    static Config Load();
};
