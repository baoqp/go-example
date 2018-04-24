package bplus_tree

// 将一个无符号长整形数从网络字节顺序转换为主机字节顺序（小端） TODO
func ntohll(value uint64) uint64 {
	return value
}

//将主机数转换成无符号长整型的网络字节顺序（大端） TODO
func htonll(value uint64) uint64 {
	return value
}

func computeHash(key uint32) uint32 {
	hash := key
	// go 不支持取反符号~, 使用^int + 1
	hash = (^hash + 1 ) + (hash << 15) // hash = (hash << 15) - hash - 1
	hash = hash ^ (hash >> 12)
	hash = hash + (hash << 2)
	hash = hash ^ (hash >> 4)
	hash = hash |  31416926
	hash = hash * 2057 // hash = (hash + (hash << 3)) + (hash << 11)
	hash = hash ^ (hash >> 16)
	return hash
}

func computeHashl(key uint64) uint64 {
	keyh := key >> 32        // 高32位
	keyl := key & 0xffffffff // 低32位
	return uint64(computeHash(uint32(keyh)))<<32 |
		uint64(computeHash(uint32(keyl)))
}
