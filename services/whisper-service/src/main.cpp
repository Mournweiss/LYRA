#include <iostream>
#include <exception>
#include "config.h"
#include "errors.h"
#include "server.cpp"

int main(int argc, char** argv) {
    try {
        Config config = Config::Load();
        RunServer(config);
    } catch (const ConfigError& e) {
        std::cerr << "Configuration error: " << e.what() << std::endl;
        return 1;
    } catch (const MinioError& e) {
        std::cerr << "MinIO error: " << e.what() << std::endl;
        return 1;
    } catch (const WhisperError& e) {
        std::cerr << "Whisper error: " << e.what() << std::endl;
        return 1;
    } catch (const TaskError& e) {
        std::cerr << "Task error: " << e.what() << std::endl;
        return 1;
    } catch (const BaseError& e) {
        std::cerr << "Service error: " << e.what() << std::endl;
        return 1;
    } catch (const std::exception& e) {
        std::cerr << "Startup error: " << e.what() << std::endl;
        throw StartupError(e.what());
    }
    return 0;
}
