#!/usr/bin/python3

import subprocess
import base64
import os
from datetime import datetime


dispatcher_cmd = ["/home/hacker/challenge/dispatcher"]
worker_cmd = ["/home/hacker/challenge/worker"]


def xor(a, b):
    return bytes(x ^ y for x, y in zip(a, b))


def dispatcher(args):
    result = subprocess.run(dispatcher_cmd + args, capture_output=True, text=True)
    return base64.b64decode(result.stdout.split()[1])


def check_pad(ct, isprint=False):
    text = f"TASK: {base64.b64encode(ct).decode()}\n"
    result = subprocess.run(worker_cmd, input=text, capture_output=True, text=True)
    if isprint and result.returncode == 0:
        print(result.stdout)
    return result.returncode == 0


def aes_pao_dfs(ct_prefix, ct1, ct2, idx, is2, check):
    if idx == -1:
        return is2

    if idx < 0 or idx > 15:
        return None

    ct1_new = ct1[:]
    for i in range(idx+1, 16):
        ct1_new[i] = is2[i] ^ (16 - idx)

    for i in range(256):
        ct1_new[idx] = i

        if not check(ct_prefix + ct1_new + ct2):
            continue

        is2 = is2[:]
        is2[idx] = ct1_new[idx] ^ (16 - idx)

        print(f"idx={idx:03}   i={i:03}   is2={is2.hex()}   pt2={xor(ct1, is2)}", flush=True)

        res = aes_pao_dfs(ct_prefix, ct1, ct2, idx - 1, is2, check)
        if res != None:
            return res
    
    return None



ct = dispatcher(["flag"])
ct = bytearray(ct)
print(f"ct: {ct.hex()}")

flag = ""

for i in range(len(ct) - 16 - 16, -1, -16):
    ct_prefix, ct1, ct2, ct_suffix = ct[:i], ct[i:i+16], ct[i+16:i+16+16], ct[i+16+16:]
    is2 = bytearray(b"\x00" * 16)

    is2 = aes_pao_dfs(ct_prefix, ct1, ct2, 16-1, is2, check_pad)
    if is2 == None:
        print("NOT FOUND")
        break

    flag = xor(ct1, is2).decode() + flag
    print(f"flag: {flag}   {flag.encode()}")
    check_pad(ct, True)
