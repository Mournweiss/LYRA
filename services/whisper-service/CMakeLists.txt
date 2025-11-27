cmake_minimum_required(VERSION 3.16)
project(whisper_grpc_service LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)
find_package(OpenSSL REQUIRED)

include_directories(${Protobuf_INCLUDE_DIRS})
include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(include)

set(PROTO_FILE ../../proto/service.proto)

add_executable(whisper_grpc_server
    src/main.cpp
    src/config.cpp
    src/handlers/task.cpp
    src/handlers/health.cpp
    src/clients/minio.cpp
    src/clients/whisper.cpp
    src/server.cpp
    /usr/local/include/INIReader.cpp
    include/service.pb.cc
    include/service.grpc.pb.cc
)

find_library(ABSL_STRINGS NAMES absl_strings)
find_library(ABSL_BASE NAMES absl_base)

find_package(unofficial-curlpp QUIET)
find_package(unofficial-inih QUIET)
find_package(nlohmann_json QUIET)
find_package(pugixml QUIET)
find_package(minio-cpp REQUIRED)

target_include_directories(whisper_grpc_server PRIVATE
    include
    /usr/local/include
    /usr/include
)

target_link_libraries(whisper_grpc_server PRIVATE
    gRPC::grpc++
    protobuf::libprotobuf
    minio-cpp::minio-cpp
    pugixml::pugixml
    unofficial::curlpp::curlpp
    unofficial::inih::inireader
    CURL::libcurl
    OpenSSL::SSL
    OpenSSL::Crypto
    /usr/local/lib/libinih.a
    ${ABSL_STRINGS}
    ${ABSL_BASE}
)
