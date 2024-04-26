package genetic_operators

// metodos a serem testados podem ter mais de um parametro de entrada
// uma seed representa uma lista de paramentros de entrada
func Crossover(seed1, seed2 []string) ([]string, []string) {
	var crossoverSeed1 []string
	var crossoverSeed2 []string

	if len(seed1) == len(seed2) {
		for i := range seed1 {
			smallest := smallestSize(seed1[i], seed2[i])
			// analisar a necessidade desta verificacao
			if smallest == 0 {
				crossoverSeed1 = seed1
				crossoverSeed2 = seed2
				//	funcionando como esperado
			} else if smallest == 1 {
				var tempSeed1 string
				var tempSeed2 string
				for j := 0; j < smallest; j++ {
					tempSeed1 = seed1[i][j:j+1] + seed2[i][j:j+1]
					tempSeed2 = seed2[i][j:j+1] + seed1[i][j:j+1]
				}
				crossoverSeed1 = append(crossoverSeed1, tempSeed1)
				crossoverSeed2 = append(crossoverSeed2, tempSeed2)
			} else {
				var tempSeed1 string
				var tempSeed2 string
				for j := 0; j < smallest; j++ {

					if j%2 == 0 {
						tempSeed1 += seed2[i][j : j+1]
						tempSeed2 += seed1[i][j : j+1]
					} else {
						tempSeed1 += seed1[i][j : j+1]
						tempSeed2 += seed2[i][j : j+1]
					}
				}
				crossoverSeed1 = append(crossoverSeed1, tempSeed1)
				crossoverSeed2 = append(crossoverSeed2, tempSeed2)
			}
		}
	}

	return crossoverSeed1, crossoverSeed2
}

func smallestSize(seed1, seed2 string) int {
	if len(seed1) >= len(seed2) {
		return len(seed2)
	}

	return len(seed1)
}
