#!/bin/bash

set -euo pipefail

# ANSI color codes
COLOR_INFO="\033[0m"       # White
COLOR_WARN="\033[1;33m"    # Yellow
COLOR_ERROR="\033[1;31m"   # Red
COLOR_SUCCESS="\033[1;32m" # Green
COLOR_RESET="\033[0m"

info()    { echo -e "${COLOR_INFO}$*${COLOR_RESET}" >&2; }
warn()    { echo -e "${COLOR_WARN}$*${COLOR_RESET}" >&2; }
error()   { echo -e "${COLOR_ERROR}$*${COLOR_RESET}" >&2; exit 1; }
success() { echo -e "${COLOR_SUCCESS}$*${COLOR_RESET}" >&2; }

install_system_dependencies() {
    info "Installing system dependencies..."

    apt-get update && \
    apt-get install -y \
        build-essential \
        cmake \
        git \
        protobuf-compiler \
        libprotobuf-dev \
        libprotoc-dev \
        libabsl-dev \
        pkg-config \
        libssl-dev \
        wget \
        libcurl4-openssl-dev \
        nlohmann-json3-dev

    success "System dependencies installed successfully"
}

build_inih_library() {
    info "Building inih library..."

    git clone https://github.com/benhoyt/inih.git /tmp/inih

    cd /tmp/inih

    gcc -c ini.c -o ini.o
    ar rcs libinih.a ini.o

    cp libinih.a /usr/local/lib/
    cp ini.h /usr/local/include/
    cp cpp/INIReader.h /usr/local/include/
    cp cpp/INIReader.cpp /usr/local/include/

    sed -i 's|../ini.h|ini.h|' /usr/local/include/INIReader.cpp

    nm /usr/local/lib/libinih.a | grep -q ini_parse || {
        error "ini_parse symbol not found in libinih.a"
        exit 1
    }

    # Cleanup
    cd /
    rm -rf /tmp/inih

    info "inih library built and installed successfully"
}

build_curlpp_library() {
    info "Building curlpp static library..."

    git clone https://github.com/jpbarrette/curlpp.git /tmp/curlpp

    cd /tmp/curlpp

    mkdir build && cd build
    cmake .. -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=OFF -DCMAKE_INSTALL_PREFIX=/usr/local
    make -j$(nproc)
    make install

    if [ ! -f "/usr/local/lib/libcurlpp.a" ]; then
        error "Static curlpp library not found at /usr/local/lib/libcurlpp.a"
        exit 1
    fi

    cd /
    rm -rf /tmp/curlpp

    info "curlpp static library built and installed successfully"
}

build_pugixml_static_library() {
    info "Building pugixml static library..."

    git clone -b v1.13 https://github.com/zeux/pugixml.git /tmp/pugixml

    cd /tmp/pugixml

    mkdir build && cd build
    cmake .. -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=OFF -DCMAKE_INSTALL_PREFIX=/usr/local
    make -j$(nproc)
    make install

    # Verify static library was created
    if [ ! -f "/usr/local/lib/libpugixml.a" ]; then
        error "Static pugixml library not found at /usr/local/lib/libpugixml.a"
        exit 1
    fi

    cd /
    rm -rf /tmp/pugixml

    info "pugixml static library built and installed successfully"
}

create_cmake_configs() {
    info "Creating CMake configuration files..."

    create_inih_cmake_config
    create_curlpp_cmake_config
    create_system_libs_cmake_config

    info "CMake configuration files created successfully"
}

create_inih_cmake_config() {
    info "Creating inih CMake configuration..."

    mkdir -p /usr/local/lib/cmake/unofficial-inih

    cat > /usr/local/lib/cmake/unofficial-inih/unofficial-inihConfig.cmake << 'EOF'
set(unofficial-inih_FOUND TRUE)
set(unofficial-inih_INCLUDE_DIRS "/usr/local/include")
set(unofficial-inih_LIBRARIES inih)
set(unofficial-inih_VERSION "1.0.0")

if(NOT TARGET unofficial::inih::inireader)
  add_library(unofficial::inih::inireader UNKNOWN IMPORTED)
  set_target_properties(unofficial::inih::inireader PROPERTIES
    IMPORTED_LOCATION "/usr/local/lib/libinih.a"
    INTERFACE_INCLUDE_DIRECTORIES "/usr/local/include"
  )
endif()
EOF
}

create_curlpp_cmake_config() {
    info "Creating curlpp CMake configuration..."

    mkdir -p /usr/local/lib/cmake/unofficial-curlpp

    cat > /usr/local/lib/cmake/unofficial-curlpp/unofficial-curlppConfig.cmake << 'EOF'
set(unofficial-curlpp_FOUND TRUE)
set(unofficial-curlpp_INCLUDE_DIRS "/usr/local/include")
set(unofficial-curlpp_LIBRARIES curlpp)
set(unofficial-curlpp_VERSION "1.0.0")

include(CMakeFindDependencyMacro)
find_dependency(CURL REQUIRED)

if(NOT TARGET unofficial::curlpp::curlpp)
  add_library(unofficial::curlpp::curlpp STATIC IMPORTED)
  set_target_properties(unofficial::curlpp::curlpp PROPERTIES
    IMPORTED_LOCATION "/usr/local/lib/libcurlpp.a"
    INTERFACE_INCLUDE_DIRECTORIES "/usr/local/include"
  )
endif()
EOF
}

create_system_libs_cmake_config() {
    info "Creating system libraries CMake configurations..."

    mkdir -p /usr/local/lib/cmake/nlohmann_json
    cat > /usr/local/lib/cmake/nlohmann_json/nlohmann_jsonConfig.cmake << 'EOF'
set(nlohmann_json_FOUND TRUE)
set(nlohmann_json_INCLUDE_DIRS "/usr/include")
set(nlohmann_json_VERSION "3.11.0")

if(NOT TARGET nlohmann_json::nlohmann_json)
  add_library(nlohmann_json::nlohmann_json INTERFACE IMPORTED)
  set_target_properties(nlohmann_json::nlohmann_json PROPERTIES
    INTERFACE_INCLUDE_DIRECTORIES "/usr/include"
  )
endif()
EOF

    mkdir -p /usr/local/lib/cmake/pugixml
    cat > /usr/local/lib/cmake/pugixml/pugixmlConfig.cmake << 'EOF'
set(pugixml_FOUND TRUE)
set(pugixml_INCLUDE_DIRS "/usr/local/include")
set(pugixml_LIBRARIES pugixml)
set(pugixml_VERSION "1.13.0")

if(NOT TARGET pugixml::pugixml)
  add_library(pugixml::pugixml STATIC IMPORTED)
  set_target_properties(pugixml::pugixml PROPERTIES
    IMPORTED_LOCATION "/usr/local/lib/libpugixml.a"
    INTERFACE_INCLUDE_DIRECTORIES "/usr/local/include"
  )
endif()
EOF

    mkdir -p /usr/local/lib/cmake/unofficial-inih-system
    cat > /usr/local/lib/cmake/unofficial-inih-system/unofficial-inihConfig.cmake << 'EOF'
set(unofficial-inih_FOUND TRUE)
set(unofficial-inih_INCLUDE_DIRS "/usr/include")
set(unofficial-inih_LIBRARIES inih)
set(unofficial-inih_VERSION "1.0.0")

if(NOT TARGET unofficial::inih::inireader)
  add_library(unofficial::inih::inireader INTERFACE IMPORTED)
  set_target_properties(unofficial::inih::inireader PROPERTIES
    INTERFACE_INCLUDE_DIRECTORIES "/usr/local/include"
  )
endif()
EOF
}

build_grpc() {
    info "Building gRPC..."

    git clone -b v1.48.0 https://github.com/grpc/grpc /grpc
    cd /grpc
    git submodule update --init
    mkdir -p cmake/build && cd cmake/build
    cmake ../.. -DgRPC_INSTALL=ON -DgRPC_BUILD_TESTS=OFF
    make -j$(nproc)
    make install

    cd /
    rm -rf /grpc

    info "gRPC built and installed successfully"
}

build_minio_cpp() {
    info "Building minio-cpp..."

    git clone https://github.com/minio/minio-cpp.git /tmp/minio-cpp
    cd /tmp/minio-cpp
    mkdir build && cd build
    cmake .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr/local
    make -j$(nproc)
    make install

    ls -la /usr/local/include/ | head -10

    create_minio_cpp_cmake_config

    cd /
    rm -rf /tmp/minio-cpp

    info "minio-cpp built and installed successfully"
}

create_minio_cpp_cmake_config() {
    info "Creating minio-cpp CMake configuration..."

    mkdir -p /usr/local/lib/cmake/minio-cpp

    cat > /usr/local/lib/cmake/minio-cpp/minio-cppConfig.cmake << 'EOF'
set(minio-cpp_FOUND TRUE)
set(minio-cpp_INCLUDE_DIRS "/usr/local/include")
set(minio-cpp_LIBRARIES minio-cpp)
set(minio-cpp_VERSION "1.0.0")

if(NOT TARGET minio-cpp::minio-cpp)
  add_library(minio-cpp::minio-cpp STATIC IMPORTED)
  set_target_properties(minio-cpp::minio-cpp PROPERTIES
    IMPORTED_LOCATION "/usr/local/lib/libminiocpp.a"
    INTERFACE_INCLUDE_DIRECTORIES "/usr/local/include"
  )
endif()
EOF
}

build_application() {
    info "Building application..."

    protoc -I./proto-context --cpp_out=./include --grpc_out=./include \
           --plugin=protoc-gen-grpc="$(which grpc_cpp_plugin)" \
           ./proto-context/service.proto

    cmake -DCMAKE_PREFIX_PATH="/usr/local/lib/cmake/grpc:/usr/local/lib/cmake:/usr/local" .
    make -j$(nproc)

    info "Application built successfully"
}
