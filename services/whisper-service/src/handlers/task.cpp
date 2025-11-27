#include "handlers/task.h"
#include <string>
#include <miniocpp/client.h>
#include "clients/minio.h"
#include "clients/whisper.h"
#include "errors.h"

bool process_transcription_task(minio::s3::Client& minio_client, const std::string& minio_bucket, const std::string& file_key, std::string& result, std::string& error) {
    std::string local_path = "/tmp/" + file_key.substr(file_key.find_last_of("/") + 1);
    try {
        if (!download_file_from_minio(minio_client, minio_bucket, file_key, local_path)) {
            throw TaskError("Failed to download file from MinIO");
        }
        result = transcribe_file_with_whisper(local_path);
        error.clear();
        std::remove(local_path.c_str());
        return true;
    } catch (const MinioError& e) {
        error = e.what();
        return false;
    } catch (const WhisperError& e) {
        error = e.what();
        return false;
    } catch (const TaskError& e) {
        error = e.what();
        return false;
    } catch (const std::exception& e) {
        error = e.what();
        return false;
    }
}
