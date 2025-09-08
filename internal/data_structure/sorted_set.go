package data_structure

type SortedSet struct {
	key      string
	skipList *Skiplist
	dict     map[string]float64 // map val to score
}

func NewSortedSet(key string) *SortedSet {
	return &SortedSet{
		key:      key,
		skipList: CreateSkiplist(),
		dict:     map[string]float64{},
	}
}

func (z *SortedSet) Add(score float64, el string) int {
	if len(el) == 0 {
		return 0
	}

	if curScore, exist := z.dict[el]; exist {
		if curScore != score {
			zNode := z.skipList.UpdateScore(curScore, el, score)
			z.dict[el] = zNode.score
		}
		return 1
	}

	zNode := z.skipList.Insert(score, el)
	z.dict[el] = zNode.score
	return 1
}

func (z *SortedSet) GetRank(el string, reversed bool) (int64, float64) {
	score, exist := z.dict[el]
	if !exist {
		return -1, 0
	}
	rank := int64(z.skipList.GetRank(score, el))
	if reversed {
		rank = int64(z.skipList.length) - 1
	} else {
		rank--
	}
	return rank, score
}

func (z *SortedSet) GetScore(el string) (int, float64) {
	score, exist := z.dict[el]
	if !exist {
		return -1, 0
	}
	return 0, score
}

func (z *SortedSet) Len() int {
	return len(z.dict)
}
