from minio import Minio
from minio.error import S3Error
import logging
from errors import MinioError

def get_minio_client(config):
    return Minio(
        config.MINIO_ENDPOINT,
        access_key=config.MINIO_ACCESS_KEY,
        secret_key=config.MINIO_SECRET_KEY,
        secure=False
    )

def upload_file(minio_client, bucket, file_key, local_path):
    
    try:
        minio_client.fput_object(
            bucket,
            file_key,
            local_path
        )
        logging.info(f"Uploaded file to MinIO: {file_key}")
        return True

    except S3Error as e:
        logging.error(f"MinIO upload failed: {e}")
        raise MinioError(f"MinIO upload failed: {e}")
