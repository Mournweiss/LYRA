import logging
from telegram import Update
from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes
import grpc
import service_pb2
import service_pb2_grpc
from config import Config
from errors import ConfigError, GRPCError

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await update.message.reply_text("Send me an audio or video file for transcription.")

async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await update.message.reply_text("Processing (stub)...")
    
    try:
        config = context.application.bot_data["config"]

        with grpc.insecure_channel(config.API_GATEWAY_ADDRESS) as channel:
            stub = service_pb2_grpc.WhisperServiceStub(channel)
            req = service_pb2.TranscribeRequest(file_content=b"", file_name="stub.wav")
            resp = stub.Transcribe.future(req)
            result = await context.application.run_in_executor(None, resp.result)
            await update.message.reply_text(f"Transcription: {result.text}")

    except grpc.RpcError as e:
        await update.message.reply_text(f"gRPC error: {e}")
        raise GRPCError(f"gRPC call failed: {e}")

    except Exception as e:
        await update.message.reply_text(f"Internal error: {e}")
        raise

def main():
    logging.basicConfig(level=logging.INFO)

    try:
        config = Config.load()

    except ConfigError as e:
        logging.error(f"Configuration error: {e}")
        exit(1)

    app = Application.builder().token(config.TELEGRAM_TOKEN).build()
    app.bot_data["config"] = config
    app.add_handler(CommandHandler("start", start))
    app.add_handler(MessageHandler(filters.ALL, handle_message))
    app.run_polling()

if __name__ == "__main__":
    main()
