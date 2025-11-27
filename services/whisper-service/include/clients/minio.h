#pragma once
#include <string>
#include <miniocpp/client.h>

bool download_file_from_minio(minio::s3::Client& minio_client, const std::string& bucket, const std::string& file_key, const std::string& local_path);
