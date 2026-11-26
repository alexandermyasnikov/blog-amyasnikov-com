type ResKey struct {
	Index int
	Key   *big.Int
}

func main() {
	log.SetFlags(log.Lmicroseconds)

	primes := factorize(N)
	log.Printf("primes: %v -> %v", N, primes)

	ch := make(chan ResKey)

	for i := range primes {
		go func(i int) {
			prime := primes[i]
			key, err := findKeyModP(msgs, encs, prime)
			log.Printf("key mod: %v - %v - %v", err, prime, key)
			ch <- ResKey{Index: i, Key: key}
		}(i)
	}

	keyMods := make([]*big.Int, len(primes))
	for range primes {
		res := <-ch
		keyMods[res.Index] = res.Key
	}

	log.Printf("primes: %v", primes)
	log.Printf("keyMods: %v", keyMods)

	key, err := crt(primes, keyMods)
	log.Printf("crt: %v - %v - %s", err, key, key.Bytes())
}