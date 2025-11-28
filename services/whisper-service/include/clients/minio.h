#pragma once
#include <string>
#include <memory>
#include <miniocpp/client.h>

std::unique_ptr<minio::s3::Client> create_minio_client(const std::string& host, int port, const std::string& access_key, const std::string& secret_key);
bool download_file_from_minio(minio::s3::Client& minio_client, const std::string& bucket, const std::string& file_key, const std::string& local_path);
