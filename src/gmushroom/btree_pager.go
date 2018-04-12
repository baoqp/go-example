package gmushroom

import "sync"

const BTreePageBucketMax = 8
type BTreePageBucket struct {
	mutex sync.Mutex

}
