#!/usr/bin/ python3

n = 5912718291679762008847883587848216166109
e = 876603837240112836821145245971528442417
print(f"n={n}")
print(f"e={e}")

p = 59644326261100157131
q = 99132954671935298039
print(f"p={p}")
print(f"q={q}")

assert n == p * q

phi = (p - 1) * (q - 1)
print(f"phi={phi}")

d = pow(e, -1, phi)
print(f"d={d}")
