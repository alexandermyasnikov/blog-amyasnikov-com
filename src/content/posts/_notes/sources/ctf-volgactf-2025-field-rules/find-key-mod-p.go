func encrypt(m, key, k, n *big.Int) *big.Int {
	e := new(big.Int).Set(m)
	for range ROUNDS - 1 {
		e.Add(e, key)
		e.Exp(e, k, n)
	}
	e.Add(e, key)
	e.Mod(e, n)
	return e
}

func findKeyModP(msgs, encs []*big.Int, p *big.Int) (*big.Int, error) {
	msgsMod := make([]*big.Int, len(msgs))
	encsMod := make([]*big.Int, len(encs))

	for i := range msgs {
		msgsMod[i] = new(big.Int).Mod(msgs[i], p)
		encsMod[i] = new(big.Int).Mod(encs[i], p)
	}

	for candidate := big.NewInt(0); candidate.Cmp(p) < 0; candidate.Add(candidate, big.NewInt(1)) {
		found := true
		for i := range msgsMod {
			e := encrypt(msgsMod[i], candidate, KEXP, p)

			if e.Cmp(encsMod[i]) != 0 {
				found = false
				if i > 0 {
					log.Printf("skip %v - %v", i, candidate)
				}
				break
			}
		}
		if found {
			return candidate, nil
		}
	}

	return nil, fmt.Errorf("no key mod p found")
}