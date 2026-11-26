#!/usr/bin/ python3
import binascii

encoded_string = open("out.txt").read().strip()
encoded_string = binascii.unhexlify(encoded_string.encode()).decode()


def decode_base_256(string):
    base = 256
    encoded_value = 0

    for char in string:
        encoded_value = encoded_value * 727 + ord(char)

    encoded_string = ""
    while encoded_value > 0:
        encoded_string = chr(encoded_value % base) + encoded_string
        encoded_value //= base

    return encoded_string


decoded_string = decode_base_256(encoded_string)
print(decoded_string)
