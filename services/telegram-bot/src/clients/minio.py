from minio import Minio
from minio.error import S3Error
import logging
import time
from errors import MinioError

def get_minio_client(config):

    try:
        client = Minio(
            config.MINIO_ENDPOINT,
            access_key=config.MINIO_ACCESS_KEY,
            secret_key=config.MINIO_SECRET_KEY,
            secure=False,
            region=config.MINIO_REGION
        )

        logging.info(f"MinIO client initialized for endpoint: {config.MINIO_ENDPOINT} (region: {config.MINIO_REGION})")

        try:
            buckets = client.list_buckets()
            logging.info(f"MinIO connection successful, found {len(buckets)} buckets")

        except Exception as e:
            logging.warning(f"MinIO connection failed: {e}")

        return client

    except Exception as e:
        logging.error(f"Failed to create MinIO client: {e}")
        raise MinioError(f"Failed to create MinIO client: {e}")

def upload_file(minio_client, bucket, file_key, local_path, max_retries=3):

    for attempt in range(max_retries):

        try:
            bucket_exists = False

            try:
                bucket_exists = minio_client.bucket_exists(bucket)
                logging.debug(f"Bucket {bucket} exists: {bucket_exists}")

            except Exception as e:
                logging.warning(f"Failed to check bucket existence (attempt {attempt + 1}): {e}")

                if attempt < max_retries - 1:
                    time.sleep(1)
                    continue

                raise MinioError(f"Failed to check bucket existence after {max_retries} attempts: {e}")

            if not bucket_exists:

                try:
                    minio_client.make_bucket(bucket)
                    logging.info(f"Created MinIO bucket: {bucket}")

                except Exception as e:
                    logging.error(f"Failed to create bucket {bucket} (attempt {attempt + 1}): {e}")

                    if attempt < max_retries - 1:
                        time.sleep(1)
                        continue

                    raise MinioError(f"Failed to create MinIO bucket {bucket} after {max_retries} attempts: {e}")

            minio_client.fput_object(
                bucket,
                file_key,
                local_path
            )

            logging.info(f"Uploaded file to MinIO: {file_key}")
            return True

        except S3Error as e:
            logging.error(f"MinIO upload failed (attempt {attempt + 1}): {e}")

            if attempt < max_retries - 1:
                time.sleep(1)
                continue

            raise MinioError(f"MinIO upload failed after {max_retries} attempts: {e}")

        except Exception as e:
            logging.error(f"Unexpected error during MinIO upload (attempt {attempt + 1}): {e}")

            if attempt < max_retries - 1:
                time.sleep(1)
                continue

            raise MinioError(f"Unexpected MinIO error after {max_retries} attempts: {e}")

    raise MinioError(f"Failed to upload file after {max_retries} attempts")
