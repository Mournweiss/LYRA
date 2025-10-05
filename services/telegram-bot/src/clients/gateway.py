import grpc
import service_pb2
import service_pb2_grpc
from errors import GatewayError, GRPCError

async def create_transcription_task(config, task_id, file_key):

    try:

        with grpc.insecure_channel(config.API_GATEWAY_ADDRESS) as channel:
            stub = service_pb2_grpc.WhisperServiceStub(channel)

            req = service_pb2.CreateTranscriptionTaskRequest(
                task_id=task_id,
                file_key=file_key
            )

            resp = stub.CreateTranscriptionTask(req)

            if resp.error:
                raise GatewayError(f"API Gateway error: {resp.error}")

            return resp
            
    except grpc.RpcError as e:
        raise GRPCError(f"gRPC call failed: {e}")
