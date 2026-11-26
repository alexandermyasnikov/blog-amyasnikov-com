#!/run/workspace/bin/python3

import os

from base64 import b64encode
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad
from Crypto.Random import get_random_bytes

key = open("/home/hacker/challenge/.key", "rb").read()
cipher = AES.new(key=key, mode=AES.MODE_CBC)
ciphertext = cipher.iv + cipher.encrypt(pad(b"sleep", cipher.block_size))

print(f"TASK: {b64encode(ciphertext).decode()}")