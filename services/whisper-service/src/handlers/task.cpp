#include "handlers/task.h"
#include <string>
#include <miniocpp/client.h>
#include "clients/minio.h"
#include "clients/whisper.h"
#include "errors.h"

bool process_transcription_task(minio::s3::Client& minio_client, const std::string& minio_bucket, const std::string& file_key, std::string& result, std::string& error) {
    std::string local_path = "/tmp/" + file_key.substr(file_key.find_last_of("/") + 1);
    std::cout << "Processing transcription task: " << file_key << " -> " << local_path << std::endl;

    try {
        std::cout << "Downloading file from MinIO..." << std::endl;
        if (!download_file_from_minio(minio_client, minio_bucket, file_key, local_path)) {
            std::cout << "Failed to download file from MinIO" << std::endl;
            throw TaskError("Failed to download file from MinIO");
        }
        std::cout << "File downloaded successfully, transcribing..." << std::endl;

        result = transcribe_file_with_whisper(local_path);
        error.clear();
        std::cout << "Transcription completed: " << result << std::endl;

        std::remove(local_path.c_str());
        std::cout << "Temporary file cleaned up" << std::endl;
        return true;
    } catch (const MinioError& e) {
        std::cerr << "MinIO error in process_transcription_task: " << e.what() << std::endl;
        error = e.what();
        return false;
    } catch (const WhisperError& e) {
        std::cerr << "Whisper error in process_transcription_task: " << e.what() << std::endl;
        error = e.what();
        return false;
    } catch (const TaskError& e) {
        std::cerr << "Task error in process_transcription_task: " << e.what() << std::endl;
        error = e.what();
        return false;
    } catch (const std::exception& e) {
        std::cerr << "Unexpected error in process_transcription_task: " << e.what() << std::endl;
        error = e.what();
        return false;
    }
}
