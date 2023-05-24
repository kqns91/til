from fastapi import FastAPI, Depends
from starlette.requests import HTTPConnection
from starlette.websockets import WebSocket
from starlette.staticfiles import StaticFiles
import asyncio
import uvloop
import uvicorn
import aioredis
import datetime

HOST = 'localhost'
PORT = 8080
REDIS_HOST = HOST
REDIS_PORT = 6379
STREAM_MAX_LEN = 1000

app = FastAPI()

app.mount('/static', StaticFiles(directory='static'), name='static')

redis = aioredis.from_url(
    'redis://localhost')

async def read_message(websocket: WebSocket, join_info: dict):
    connected = True
    is_first = True
    stream_id = '$'
    while connected:
        try:
            count = 1 if is_first else 100
            results = await redis.xread(
                streams={join_info['room']: stream_id},
                count=count,
                block=100000
            )
            for room, events in results:
                if join_info['room'] != room.decode('utf-8'):
                    continue
                for e_id, e in events:
                    now = datetime.datetime.now()
                    await websocket.send_text(f"{now.strftime('%H:%M')} {e[b'msg'].decode('utf-8')}")
                    stream_id = e_id
                if is_first:
                    is_first = False
        except:
            await redis.close()
            connected = False

async def write_message(websocket: WebSocket, join_info: dict):
    await notify(join_info, 'joined')

    connected = True
    while connected:
        try:
            data = await websocket.receive_text()
            await redis.xadd(join_info['room'],
                             {
                                 'username': join_info['username'],
                                 'msg': data
                             },
                             id=b'*',
                             maxlen=STREAM_MAX_LEN)
        except:
            await notify(join_info, 'left')
            await redis.close()
            connected = False

async def notify(join_info: dict, action: str):
    await redis.xadd(join_info['room'],
                     {'msg': f"{join_info['username']} has {action}"},
                     id=b'*',
                     maxlen=STREAM_MAX_LEN)

async def get_joininfo(username: str = None, room: str = None):
    return {"username": username, "room": room}

@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket, join_info: dict = Depends(get_joininfo)):
    await websocket.accept()
    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
    await asyncio.gather(write_message(websocket, join_info), read_message(websocket, join_info))

if __name__ == '__main__':
    uvicorn.run(app, host=HOST, port=PORT)