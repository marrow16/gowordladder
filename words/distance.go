package words

type WordDistanceMap map[string]int

func NewWordDistanceMap(word *Word, maximumLadderLength *int) WordDistanceMap {
	result := make(WordDistanceMap)
	result[word.String()] = 1
	maxDistance := word.MaxSteps()
	if maximumLadderLength != nil {
		maxDistance = *maximumLadderLength
	}
	q := make([]*Word, 0, 1024)
	head := 0
	q = append(q, word)
	for head < len(q) {
		nextWord := q[head]
		head++
		distance := result[nextWord.String()] + 1
		if distance <= maxDistance {
			for _, linkedWord := range nextWord.LinkedWords() {
				if _, ok := result[linkedWord.String()]; !ok {
					q = append(q, linkedWord)
					result[linkedWord.String()] = distance
				}
			}
		}
	}
	return result
}

func (m WordDistanceMap) Distance(word *Word) (dist int, ok bool) {
	dist, ok = m[word.String()]
	return
}

func (m WordDistanceMap) Reachable(word *Word, maximumLadderLength int) bool {
	if distance, ok := m[word.String()]; ok {
		return distance <= maximumLadderLength
	}
	return false
}

func (m WordDistanceMap) WordsAt(ladderLength int) (result []string) {
	for k, v := range m {
		if v == ladderLength {
			result = append(result, k)
		}
	}
	return result
}

func (m WordDistanceMap) Words() (result []string) {
	for k, v := range m {
		if v > 1 {
			result = append(result, k)
		}
	}
	return result
}

func (m WordDistanceMap) MaxDistance() (result int) {
	for _, v := range m {
		if v > result {
			result = v
		}
	}
	return result
}
