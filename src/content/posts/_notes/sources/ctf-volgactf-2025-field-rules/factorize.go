func factorize(n *big.Int) []*big.Int {
	num := new(big.Int).Set(n)
	divisor := big.NewInt(2)
	remainder := new(big.Int)
	tmp := new(big.Int)

	primes := make([]*big.Int, 0)

	for {
		tmp.QuoRem(num, divisor, remainder)

		if remainder.BitLen() == 0 {
			primes = append(primes, divisor)
			num.Div(num, divisor)
			divisor = big.NewInt(2)

			if num.Cmp(big.NewInt(1)) == 0 {
				break
			}
		}

		divisor.Add(divisor, big.NewInt(1))
	}

	return primes
}