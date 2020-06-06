package util

import (
	"fmt"
	"sync"
)

// 组合（通用数学方法）
// 参考：https://github.com/mxschmitt/golang-combinations/blob/master/combinations.go
// 		https://stackoverflow.com/questions/56103775/how-to-print-formatted-string-to-the-same-line-in-stdout-with-go
func Combination(set []string, size int, threads int) (subsets [][]string) {
	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	length := uint(len(set))
	total := 1 << length
	num := total / threads
	if total%threads != 0 {
		threads++
	}
	res := make([][][]string, threads)

	wg := sync.WaitGroup{}
	wg.Add(threads)

	for th := 0; th < threads; th++ {
		from := th * num
		if from == 0 {
			from = 1
		}
		end := (th+1)*num
		if end > total {
			end = total
		}
		go func(th, from, end int) {
			fmt.Printf("[#%03d] [%d, %d) (%d)\n", th, from, end, end - from)

			for bit := from; bit < end; bit++ {
				// fmt.Printf("[%d] %0.2f%%\r", th, float64((bit-from)*100)/float64(end-from))
				var subset []string
				for object := uint(0); object < length; object++ {
					// checks if object is contained in subset
					// by checking if bit 'object' is set in subsetBits
					if (bit>>object)&1 == 1 {
						// add object to subset
						subset = append(subset, set[object])
					}
				}
				// add subset to subsets
				if len(subset) == size {
					res[th] = append(res[th], subset)
				}
			}
			wg.Done()
		}(th, from, end)
	}

	wg.Wait()

	for _, r := range res {
		subsets = append(subsets, r...)
	}
	return subsets
}
