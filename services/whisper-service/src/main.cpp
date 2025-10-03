#include <iostream>
#include <memory>
#include <string>

#include <grpcpp/grpcpp.h>
#include "service.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using lyra::WhisperService;
using lyra::TranscribeRequest;
using lyra::TranscribeResponse;
using lyra::HealthCheckRequest;
using lyra::HealthCheckResponse;

class WhisperServiceImpl final : public WhisperService::Service {
    Status Transcribe(ServerContext* context, const TranscribeRequest* request, TranscribeResponse* response) override {
        response->set_text("Test transcription result");
        response->set_error("");
        return Status::OK;
    }

    Status HealthCheck(ServerContext* context, const HealthCheckRequest* request, HealthCheckResponse* response) override {
        response->set_status("SERVING");
        return Status::OK;
    }
};

void RunServer() {
    std::string server_address("0.0.0.0:50052");
    WhisperServiceImpl service;

    ServerBuilder builder;
    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);

    std::unique_ptr<Server> server(builder.BuildAndStart());
    std::cout << "WhisperService gRPC server listening on " << server_address << std::endl;
    server->Wait();
}

int main(int argc, char** argv) {
    RunServer();
    return 0;
}
