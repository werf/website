#!/usr/bin/env python
import asyncio
import time
import socketio
import os
import urllib
import random

loop = asyncio.get_event_loop()
sio = socketio.AsyncClient(logger=True)
start_timer = None

LOGIN = os.environ.get('BOT_LOGIN', 'bash')
PASSWORD = os.environ.get('BOT_PASSWORD', 'bash')
SIO_URL = os.environ.get('SIO_URL', 'http://localhost:5000')
cookie_jar= None

def get_login(url, login, password):
    data = urllib.parse.urlencode({'user_email': login, 'user_password': password})
    data = data.encode('ascii')
    req = urllib.request.Request(url, data)
    cookie = None
    try:
        response = urllib.request.urlopen(req)
        cookie = response.info()['Set-Cookie']
        content_type = response.info()['Content-Type']
    except urllib.error.HTTPError as err:
        print("err status: {0}".format(err))
    return cookie

async def send_responce(responce):
    await sio.emit("chat_message", responce)

@sio.event
async def is_online(data):
    print('I received a message!')
    print(data)
    hello_list = ["HELLO!!","ALLOHA!!","Namaste :D"]
    await send_responce(random.choice(hello_list))

@sio.event
async def connect():
    print('connected to server')

async def start_server():
    global SIO_URL
    global cookie_jar
    auth_headers = { 'cookie' : cookie_jar }
    await sio.connect(SIO_URL, headers = auth_headers)
    await sio.wait()

if __name__ == '__main__':
    cookie_jar = get_login(SIO_URL + '/login', LOGIN, PASSWORD).split(';')[0]
    print(cookie_jar)
    loop.run_until_complete(start_server())