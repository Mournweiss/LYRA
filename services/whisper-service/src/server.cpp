#include <iostream>
#include <memory>
#include <string>
#include <grpcpp/grpcpp.h>
#include <stdexcept>
#include "service.grpc.pb.h"
#include "config.h"
#include "errors.h"
#include "clients/minio.h"
#include "clients/whisper.h"
#include "handlers/task.h"
#include "handlers/health.h"

using grpc::Server;
using grpc::ServerBuilder;
using lyra::WhisperService;
using lyra::TranscribeRequest;
using lyra::TranscribeResponse;
using lyra::HealthCheckRequest;
using lyra::HealthCheckResponse;

class WhisperServiceImpl final : public WhisperService::Service {
private:
    std::unique_ptr<minio::s3::Client> minio_client_;
    std::string minio_bucket_;

public:
    WhisperServiceImpl(const Config& config)
        : minio_client_(nullptr),
          minio_bucket_(config.minio_bucket) {
        try {
            int port = std::stoi(config.minio_port);
            minio_client_ = create_minio_client(config.minio_host, port, config.minio_access_key, config.minio_secret_key);
        } catch (const std::exception& e) {
            std::cerr << "Failed to initialize WhisperServiceImpl: " << e.what() << std::endl;
            throw;
        }
    }

    grpc::Status Transcribe(grpc::ServerContext* context, const TranscribeRequest* request, TranscribeResponse* response) override {
        std::cout << "Received Transcribe gRPC request" << std::endl;
        try {
            std::string file_key = request->file_key();
            std::cout << "File key: " << file_key << std::endl;

            if (file_key.empty()) {
                std::cout << "Error: file_key is empty" << std::endl;
                response->set_error("file_key is required");
                return grpc::Status(grpc::INVALID_ARGUMENT, "file_key is required");
            }

            if (!minio_client_) {
                std::cout << "Error: MinIO client not initialized" << std::endl;
                response->set_error("Internal server error: MinIO client not initialized");
                return grpc::Status(grpc::INTERNAL, "MinIO client not initialized");
            }

            std::cout << "Starting transcription task processing..." << std::endl;
            std::string result, error;
            bool ok = process_transcription_task(*minio_client_, minio_bucket_, file_key, result, error);
            if (!ok) {
                response->set_error(error);
                return grpc::Status(grpc::INTERNAL, error);
            }
            response->set_text(result);
            response->set_error("");
            return grpc::Status::OK;
        } catch (const MinioError& e) {
            response->set_error(std::string("MinIO error: ") + e.what());
            return grpc::Status(grpc::INTERNAL, e.what());
        } catch (const WhisperError& e) {
            response->set_error(std::string("Whisper error: ") + e.what());
            return grpc::Status(grpc::INTERNAL, e.what());
        } catch (const TaskError& e) {
            response->set_error(std::string("Task error: ") + e.what());
            return grpc::Status(grpc::INTERNAL, e.what());
        } catch (const std::exception& e) {
            response->set_error(std::string("Exception: ") + e.what());
            return grpc::Status(grpc::INTERNAL, e.what());
        }
    }

    grpc::Status HealthCheck(grpc::ServerContext* context, const HealthCheckRequest* request, HealthCheckResponse* response) override {
        response->set_status(handle_health_check());
        return grpc::Status::OK;
    }
};

void RunServer(const Config& config) {
    try {
        std::string server_address = "0.0.0.0:" + config.service_port;
        WhisperServiceImpl service(config);
        ServerBuilder builder;
        builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
        builder.RegisterService(&service);
        std::unique_ptr<Server> server(builder.BuildAndStart());
        if (!server) {
            throw StartupError("Failed to start gRPC server");
        }
        std::cout << "WhisperService gRPC server listening on " << server_address << std::endl;
        server->Wait();
    } catch (const std::exception& e) {
        std::cerr << "Server startup error: " << e.what() << std::endl;
        throw;
    }
}
