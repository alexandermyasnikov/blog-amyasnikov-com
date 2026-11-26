func crt(p, k []*big.Int) (*big.Int, error) {
	N := new(big.Int).SetInt64(1)
	for _, v := range p {
		N.Mul(N, v)
	}

	x := new(big.Int)
	for i := range p {
		Ni := new(big.Int).Div(N, p[i])
		inv := new(big.Int)
		if inv.ModInverse(Ni, p[i]) == nil {
			return nil, fmt.Errorf("no inverse for %v mod %v", Ni, p[i])
		}
		x.Add(x, new(big.Int).Mul(k[i], new(big.Int).Mul(Ni, inv)))
	}

	return x.Mod(x, N), nil
}