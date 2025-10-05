#include <iostream>
#include <memory>
#include <string>
#include <grpcpp/grpcpp.h>
#include "service.grpc.pb.h"
#include "config.h"
#include "errors.h"
#include <minio/minio.hpp>
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
public:
    WhisperServiceImpl(const Config& config)
        : minio_client_(minio::s3::Client(
            config.minio_host + ":" + config.minio_port,
            config.minio_access_key,
            config.minio_secret_key,
            false)),
          minio_bucket_(config.minio_bucket) {}

    grpc::Status Transcribe(grpc::ServerContext* context, const TranscribeRequest* request, TranscribeResponse* response) override {
        try {
            std::string file_key = request->file_key();
            if (file_key.empty()) {
                response->set_error("file_key is required");
                return grpc::Status(grpc::INVALID_ARGUMENT, "file_key is required");
            }
            std::string result, error;
            bool ok = process_transcription_task(minio_client_, minio_bucket_, file_key, result, error);
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
private:
    minio::s3::Client minio_client_;
    std::string minio_bucket_;
};

void RunServer(const Config& config) {
    std::string server_address = "0.0.0.0:" + config.service_port;
    WhisperServiceImpl service(config);
    ServerBuilder builder;
    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);
    std::unique_ptr<Server> server(builder.BuildAndStart());
    std::cout << "WhisperService gRPC server listening on " << server_address << std::endl;
    server->Wait();
}
