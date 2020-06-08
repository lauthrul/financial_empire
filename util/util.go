package util

import (
	"sync"

	"financial_empire/logger"
)

// 组合（通用数学方法）
// 参考：https://github.com/mxschmitt/golang-combinations/blob/master/combinations.go
// 		https://stackoverflow.com/questions/56103775/how-to-print-formatted-string-to-the-same-line-in-stdout-with-go
func Combination(set []string, size uint64, threads uint64, subsets *[][]string) {
	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	length := uint64(len(set))
	total := uint64(1 << length)
	num := total / threads
	if total%threads != 0 {
		threads++
	}
	res := make([][]uint64, threads)

	wg := sync.WaitGroup{}
	wg.Add(int(threads))

	for th := uint64(0); th < threads; th++ {
		from := th * num
		if from == 0 {
			from = 1
		}
		end := (th + 1) * num
		if end > total {
			end = total
		}
		go func(th, from, end uint64) {
			logger.Logger.Printf("[#%03d] [%d, %d) (%d)\n", th, from, end, end-from)

			for bit := from; bit < end; bit++ {
				// logger.Logger.Printf("[%d] %0.2f%%\r", th, float64((bit-from)*100)/float64(end-from))
				count := uint64(0)
				for object := uint64(0); object < length; object++ {
					// checks if object is contained in subset
					// by checking if bit 'object' is set in subsetBits
					if (bit>>object)&1 == 1 {
						count++
					}
				}
				if count != size {
					continue
				}

				// add subset to subsets
				res[th] = append(res[th], bit)
			}
			wg.Done()
		}(th, from, end)
	}

	wg.Wait()

	for _, r := range res {
		for _, v := range r {
			// add object to subset
			var subset []string
			for object := uint64(0); object < length; object++ {
				if (v>>object)&1 == 1 {
					subset = append(subset, set[object])
				}
			}
			*subsets = append(*subsets, subset)
		}
	}
}
