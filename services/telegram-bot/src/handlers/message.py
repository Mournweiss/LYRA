import os
import logging
import time
from telegram import Update
from telegram.ext import ContextTypes
from clients.minio import get_minio_client, upload_file
from utils import generate_task_id
from config import Config
from .task import send_create_task_request
from errors import MinioError, GatewayError, ValidationError

async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE):
    config = context.application.bot_data["config"]
    minio_client = get_minio_client(config)

    try:
        file = None
        original_filename = None

        if update.message.document:
            file = await update.message.document.get_file()
            original_filename = update.message.document.file_name

        elif update.message.audio:
            file = await update.message.audio.get_file()
            original_filename = update.message.audio.file_name

        elif update.message.voice:
            file = await update.message.voice.get_file()
            original_filename = f"voice_{int(time.time())}.ogg"

        elif update.message.video:
            file = await update.message.video.get_file()
            original_filename = update.message.video.file_name

        else:
            raise ValidationError("Please send an audio, video, voice message, or document file.")

        local_path = f"/tmp/{original_filename}"
        await file.download_to_drive(local_path)
        task_id = generate_task_id()
        file_key = f"{task_id}/{original_filename}"
        upload_file(minio_client, config.MINIO_BUCKET, file_key, local_path)
        error = await send_create_task_request(config, task_id, file_key, update)

        if not error:
            await update.message.reply_text(f"Task created! Task ID: {task_id}\nYou can check status later.")

        try:
            os.remove(local_path)

        except Exception:
            pass

    except ValidationError as e:
        await update.message.reply_text(str(e))

    except MinioError as e:
        await update.message.reply_text(f"File upload error: {e}")

    except GatewayError as e:
        await update.message.reply_text(f"API Gateway error: {e}")

    except Exception as e:
        await update.message.reply_text(f"Internal error: {e}")
        logging.exception(e)
