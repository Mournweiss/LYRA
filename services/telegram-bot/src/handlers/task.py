from clients.gateway import create_transcription_task
from errors import GatewayError, GRPCError

async def send_create_task_request(config, task_id, file_key, update):

    try:
        resp = await create_transcription_task(config, task_id, file_key)
        return None

    except GatewayError as e:
        await update.message.reply_text(f"API Gateway error: {e}")
        return str(e)

    except GRPCError as e:
        await update.message.reply_text(f"gRPC error: {e}")
        return str(e)
        
    except Exception as e:
        await update.message.reply_text(f"Internal error: {e}")
        return str(e)
