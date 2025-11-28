#include <iostream>
#include <exception>
#include "config.h"
#include "errors.h"

void RunServer(const Config& config);

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
    } catch (const StartupError& e) {
        std::cerr << "Startup error: " << e.what() << std::endl;
        return 1;
    } catch (const BaseError& e) {
        std::cerr << "Service error: " << e.what() << std::endl;
        return 1;
    } catch (const std::exception& e) {
        std::cerr << "Unexpected error: " << e.what() << std::endl;
        return 1;
    }
    return 0;
}
