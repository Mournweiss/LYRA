#pragma once
#include <string>
#include <minio/minio.hpp>

bool process_transcription_task(minio::s3::Client& minio_client, const std::string& minio_bucket, const std::string& file_key, std::string& result, std::string& error);
