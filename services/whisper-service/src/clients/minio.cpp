#include "clients/minio.h"
#include <miniocpp/client.h>
#include <fstream>
#include <string>
#include <iostream>
#include "errors.h"

std::unique_ptr<minio::s3::Client> create_minio_client(const std::string& host, int port, const std::string& access_key, const std::string& secret_key) {
    try {
        minio::s3::BaseUrl base_url(host + ":" + std::to_string(port), false);
        auto provider = std::make_unique<minio::creds::StaticProvider>(access_key, secret_key);
        auto client = std::make_unique<minio::s3::Client>(base_url, provider.get());

        std::cout << "MinIO client initialized for " << host << ":" << port << std::endl;
        return client;
    } catch (const std::exception& e) {
        std::cerr << "Failed to initialize MinIO client: " << e.what() << std::endl;
        throw MinioError(std::string("Failed to create MinIO client: ") + e.what());
    }
}

bool download_file_from_minio(minio::s3::Client& minio_client, const std::string& bucket, const std::string& file_key, const std::string& local_path) {
    std::cout << "Starting download_file_from_minio function" << std::endl;
    std::cout << "Parameters: bucket=" << bucket << ", key=" << file_key << ", path=" << local_path << std::endl;

    if (bucket.empty()) {
        throw MinioError("Bucket name cannot be empty");
    }
    if (file_key.empty()) {
        throw MinioError("File key cannot be empty");
    }
    if (local_path.empty()) {
        throw MinioError("Local path cannot be empty");
    }

    try {
        std::cout << "Checking bucket existence..." << std::endl;

        minio::s3::BucketExistsArgs bucket_args;
        bucket_args.bucket = bucket;
        minio::s3::BucketExistsResponse bucket_resp = minio_client.BucketExists(bucket_args);

        if (!bucket_resp) {
            std::string error_msg = "Bucket does not exist: " + bucket_resp.Error().String();
            std::cerr << error_msg << std::endl;
            throw MinioError(error_msg);
        }

        std::cout << "Bucket exists, proceeding with download..." << std::endl;
        std::cout << "Creating DownloadObjectArgs..." << std::endl;
        minio::s3::DownloadObjectArgs args;
        args.bucket = bucket;
        args.object = file_key;
        args.filename = local_path;

        std::cout << "Calling DownloadObject..." << std::endl;
        minio::s3::DownloadObjectResponse resp = minio_client.DownloadObject(args);
        std::cout << "DownloadObject returned, checking response..." << std::endl;

        if (!resp) {
            std::string error_msg = "Failed to download file from MinIO: " + resp.Error().String();
            std::cerr << error_msg << std::endl;
            throw MinioError(error_msg);
        }

        std::cout << "Download successful, checking file existence..." << std::endl;
        std::ifstream file_check(local_path, std::ios::binary | std::ios::ate);

        if (!file_check.good()) {
            std::string error_msg = "Downloaded file does not exist or is not readable: " + local_path;
            std::cerr << error_msg << std::endl;
            throw MinioError(error_msg);
        }

        std::streamsize file_size = file_check.tellg();
        file_check.close();

        std::cout << "File downloaded successfully, size: " << file_size << " bytes" << std::endl;
        return true;

    } catch (const MinioError& e) {
        std::cerr << "MinIO error in download_file_from_minio: " << e.what() << std::endl;
        throw;
    } catch (const std::exception& e) {
        std::string error_msg = std::string("Unexpected exception in download_file_from_minio: ") + e.what();
        std::cerr << error_msg << std::endl;
        throw MinioError(error_msg);
    } catch (...) {
        std::string error_msg = "Unknown exception in download_file_from_minio";
        std::cerr << error_msg << std::endl;
        throw MinioError(error_msg);
    }
}
