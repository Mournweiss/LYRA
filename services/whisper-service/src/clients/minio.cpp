#include "clients/minio.h"
#include <minio/minio.hpp>
#include <fstream>
#include <string>
#include <iostream>
#include "errors.h"

bool download_file_from_minio(minio::s3::Client& minio_client, const std::string& bucket, const std::string& file_key, const std::string& local_path) {
    try {
        minio::s3::GetObjectArgs args = minio::s3::GetObjectArgs().bucket(bucket).object(file_key);
        std::ofstream ofs(local_path, std::ios::binary);
        minio::s3::GetObjectResponse get_resp = minio_client.GetObject(args, [&ofs](const char* data, size_t size) {
            ofs.write(data, size);
        });
        ofs.close();
        if (!get_resp) {
            throw MinioError("Failed to download file from MinIO: " + get_resp.Error().String());
        }
        return true;
    } catch (const std::exception& e) {
        throw MinioError(std::string("MinIO exception: ") + e.what());
    }
}
