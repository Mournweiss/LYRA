import logging
from telegram import Update
from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes
import grpc
import service_pb2
import service_pb2_grpc

TELEGRAM_TOKEN = "your-telegram-bot-token"  # TODO: load from env
API_GATEWAY_ADDRESS = "api-gateway:50051"

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await update.message.reply_text("Send me an audio or video file for transcription.")

async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await update.message.reply_text("Processing (stub)...")
    try:
        with grpc.insecure_channel(API_GATEWAY_ADDRESS) as channel:
            stub = service_pb2_grpc.WhisperServiceStub(channel)
            req = service_pb2.TranscribeRequest(file_content=b"", file_name="stub.wav")
            resp = stub.Transcribe.future(req)
            result = await context.application.run_in_executor(None, resp.result)
            await update.message.reply_text(f"Transcription: {result.text}")
    except Exception as e:
        await update.message.reply_text(f"gRPC error: {e}")

def main():
    logging.basicConfig(level=logging.INFO)
    app = Application.builder().token(TELEGRAM_TOKEN).build()
    app.add_handler(CommandHandler("start", start))
    app.add_handler(MessageHandler(filters.ALL, handle_message))
    app.run_polling()

if __name__ == "__main__":
    main()
