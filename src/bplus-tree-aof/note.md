# B+ tree

## 特性
- internal node 只有放子节点的索引，数据都在 leaf node
- leaf node在相同的深度，互相连接
- internal node的key都会出现在leaf node(B tree不会)

## 为什么要使用 B+ tree
- 因为internal node 没有放数据，所以一个内存页可以放更多的key和索引，减少cache miss
- leaf node 组成一个linked list，做范围时比较方便


## 参考资料
- [concurrent B+ tree]()https://hackmd.io/s/SyjKs-mxg)


## 源码分析

## 概述
这份代码来自于https://github.com/embedded2016/bplus-tree，此处把C代码翻译成go的实现，
此处的B+树以append only file的形式实现，主要思路是在节点的任何一次修改之后都会把节点序列化并
append到文件末尾，同时更新在父节点中的信息，由于父节点的数据也改变了，所以也需要重新序列化后写入文件，
如此迭代直到树根。所以这种方法是很浪费空间的，需要定期执行compact操作。

### 初始化
入口在open()方法中，该方法中会打开数据文件，然后从文件末尾开始每32个byte遍历根据hash值来查找Header，
Header中的offset和congfig就是B+树的根page（这里以page来代表B+树上的节点）在文件中的位置，以及该
page的大小，然后把该page载入内存。如果是首次打开数据文件，那么会实例化Header，以及新建一个Page
作为根，具体代码见treeReadHead()方法和treeWriteHead方法。

### CURD
插入和更新操作调用了insert方法，如果是更新的话会先删除老数据再插入新数据。插入操作和B+树的标准插入操作，
从根page开始，先查找在哪个子page然后载入内存，接着查找和载入，直到叶子page。待插入的数据和B+树的节点page
是写入相同的文件的，需要记录每个page数据或具体数据的大小和在文件中offset。
删除操作没有遵循标准的B+树删除操作，删除节点后没有进行树的平衡操作，而是当一个节点没有元素了直接从父节点中删除
该节点的信息。
为了支持并发，采用了读写锁进行加锁。

### compact
由于采用了AOF的方式， 每次对数据的操作不会删除截数据，因此会造成空间的浪费，需要进行compact操作以减少文件的大小。
compact操作只需简单地把整棵树遍历一遍，依次把每个节点或具体数据写入新的文件即可，且在compact时只能进行读操作。




https://github.com/begeekmyfriend/bplustree
https://github.com/NicolasLM/bplustree
https://github.com/timtadh/fs2
https://github.com/malbrain/Btree-source-code


[how the append-only btree works](http://www.bzero.se/ldapd/btree.html)