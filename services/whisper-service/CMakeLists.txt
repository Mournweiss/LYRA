cmake_minimum_required(VERSION 3.16)
project(whisper_grpc_service LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

# gRPC and Protobuf
find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)

include_directories(${Protobuf_INCLUDE_DIRS})
include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(include)

# Proto file
set(PROTO_FILE ../../proto/service.proto)

# Generate sources
add_executable(whisper_grpc_server
    src/main.cpp
    src/config.cpp # Added to resolve undefined reference to Config::Load()
    include/service.pb.cc
    include/service.grpc.pb.cc
)

find_library(ABSL_STRINGS NAMES absl_strings)
find_library(ABSL_BASE NAMES absl_base)

target_link_libraries(whisper_grpc_server
    gRPC::grpc++
    protobuf::libprotobuf
    ${ABSL_STRINGS}
    ${ABSL_BASE}
)
