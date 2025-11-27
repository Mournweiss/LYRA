#include "clients/minio.h"
#include <miniocpp/client.h>
#include <fstream>
#include <string>
#include <iostream>
#include "errors.h"

bool download_file_from_minio(minio::s3::Client& minio_client, const std::string& bucket, const std::string& file_key, const std::string& local_path) {
    try {
        minio::s3::DownloadObjectArgs args;
        args.bucket = bucket;
        args.object = file_key;
        args.filename = local_path;

        minio::s3::DownloadObjectResponse resp = minio_client.DownloadObject(args);
        if (!resp) {
            throw MinioError("Failed to download file from MinIO: " + resp.Error().String());
        }
        return true;
    } catch (const std::exception& e) {
        throw MinioError(std::string("MinIO exception: ") + e.what());
    }
}
