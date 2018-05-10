
## etcd raft库的使用

### [etcd raft library设计原理和使用](http://www.cnblogs.com/foxmailed/p/7137431.html)

etcd提供的raft库只负责把消息保存在内存中，需要自行实现网络传输层以完成消息的发送和接收

raft提供了Ready结构体，其中封装了只读的entry和message, 可以持久化到存储或提交或发送给其他peer，具体包括以下部分：
- pb.HardState：包含当前节点见过的最大的term，以及在该term的投票，当前的commit index
- Messages: 需要广播给所有peers的消息
- CommittedEntries：已经commit但还未apply到状态机的日志
- Snapshot：需要持久化的快照
使用时从node结构体提供的Ready()方法获取到ready channel，并从中获取Ready进行处理。
```go
type Ready struct {
    // The current volatile state of a Node.
    // SoftState will be nil if there is no update.
    // It is not required to consume or store SoftState.
    *SoftState

    // The current state of a Node to be saved to stable storage BEFORE
    // Messages are sent.
    // HardState will be equal to empty state if there is no update.
    pb.HardState

    // ReadStates can be used for node to serve linearizable read requests locally
    // when its applied index is greater than the index in ReadState.
    // Note that the readState will be returned when raft receives msgReadIndex.
    // The returned is only valid for the request that requested to read.
    ReadStates []ReadState

    // Entries specifies entries to be saved to stable storage BEFORE
    // Messages are sent.
    Entries []pb.Entry

    // Snapshot specifies the snapshot to be saved to stable storage.
    Snapshot pb.Snapshot

    // CommittedEntries specifies entries to be committed to a
    // store/state-machine. These have previously been committed to stable
    // store.
    CommittedEntries []pb.Entry

    // Messages specifies outbound messages to be sent AFTER Entries are
    // committed to stable storage.
    // If it contains a MsgSnap message, the application MUST report back to raft
    // when the snapshot has been received or has failed by calling ReportSnapshot.
    Messages []pb.Message

    // MustSync indicates whether the HardState and Entries must be synchronously
    // written to disk or if an asynchronous write is permissible.
    MustSync bool
}
```

应用需要对Ready的处理包括:
- 将HardState, Entries, Snapshot持久化到storage。
- 将Messages非阻塞地广播给其他peers
- 将CommittedEntries应用到状态机。
- 如果发现CommittedEntries中有成员变更类型的entry，调用node的ApplyConfChange()进行通知
- 调用Node.Advance()告诉raft node，这批状态更新处理完了，可以处理下一批更新了。

应用通过raft.StartNode()来启动raft中的一个副本，函数内部通过启动一个goroutine运行run方法来启动服务。
```go
func (n *node) run(r *raft)
```


通过调用Propose方法向raft提交请求
```go
func (n *node) Propose(ctx context.Context, data []byte) error
```


增删节点通过调用
```go
func (n *node) ProposeConfChange(ctx context.Context, cc pb.ConfChange) error
```

node结构体包含几个重要的channel:
```go
// node is the canonical implementation of the Node interface
type node struct {
    propc      chan pb.Message
    recvc      chan pb.Message
    confc      chan pb.ConfChange
    confstatec chan pb.ConfState
    readyc     chan Ready
    advancec   chan struct{}
    tickc      chan struct{}
    done       chan struct{}
    stop       chan struct{}
    status     chan chan Status

    logger Logger
}
```

- propc: propc是一个没有buffer的channel，应用通过Propose方法写入的请求被封装成Message然后push到propc中，
node的run方法从propc中取出Message，append到自己的raft log中，并且将Message放入mailbox中(raft结构体中的msgs []pb.Message)，
这个msgs会被封装在Ready中，被应用从readyc中取出来，然后通过应用自定义的transport发送出去。

- recvc: 应用自定义的transport在收到Message后需要调用
```go
func (n *node) Step(ctx context.Context, m pb.Message) error
```
把Message放入recvc中，经过一些处理后，会把需要发送的Message放入到对应peers的mailbox中。

- readyc/advancec: readyc和advancec都是没有buffer的channel，node.run()内部把相关的一些状态更新打包成Ready结构体
放入readyc中。应用从readyc中取出对相应的状态进行处理，处理完成后，调用
```go
rc.node.Advance()
```
往advancec中push一个空结构体告诉raft处理完毕，node.run()内部从advancec中得到通知后，对内部一些状态进行处理，
比如把已经持久化到storage中的entries从内存(对应type unstable struct)中删除等。

- tickc:应用定期往tickc中push空结构体，node.run()会调用tick()函数。对于leader来说，tick()会给其他peers发心跳，
对于follower来说，会检查是否需要发起选主操作。

- confc/confstatec:应用从Ready中拿出CommittedEntries，检查其如果含有成员变更类型的日志，则需要调用
```go
func (n *node) ApplyConfChange(cc pb.ConfChange) *pb.ConfState
```
把ConfChange发送到confc中，confc同样是个无buffer的channel，node.run()内部会从confc中拿出ConfChange，
然后进行真正的增减peers操作，之后将最新的成员组push到confstatec中，而ApplyConfChange函数从confstatec
取出最新的成员组返回给应用。



### [etcd-raft使用分析](http://www.opscoder.info/ectd-raft-example.html)




## etcd raft库的原理

https://godoc.org/github.com/coreos/etcd/raft

### 参考资料
知乎专栏 https://zhuanlan.zhihu.com/distributed-storage
[Etcd Raft Libary 源码阅读：Core]https://blog.neverchanje.com/2017/01/30/etcd_raft_core/

### WAL
[package wal](https://godoc.org/github.com/coreos/etcd/wal)
[etcd-raft日志管理](https://zhuanlan.zhihu.com/p/29692778)

日志项会被存储在三个地方，按照其出现的顺序分别为：unstable, WAL 和 storage
unstable: 维护协议层的日志项，保存在内存中，不稳定
WAL(write ahead log): 负责日志项的持久化存储

> 只考虑log entries的话，unstable是未落盘的，WAL是已落盘entries，storage是访问已落盘数据的interface，
具体实现上，一般是WAL加某种cache的实现。etcd自带的memoryStorage实现这个storage接口，但比较简单，是没有被
compact掉的已落盘entries在内存的一份拷贝，和传统意义cache不同，因为它有已落盘未compact掉的所有数据。
unstable不是复制数据的来源，在有follower落后、刚重启、新join的情况下，给这类follower的数据多数来自已落盘部分。
cockroachdb使用一个基于llrb的LRU cache来替代memoryStorage这个东西，WAL部分是rocksdb。

raft storage 内存中维护了那些已经被写入WAL但是未compact的日志项，同时还记录了最近一次的snapshot信息。
节点每次完成snapshot后，便可以回收该snapshot之前的所有日志项，以释放日志项占用的内存。

unstable log中的日志项来源主要有二：于Leader节点，日志项是来自客户端的更新请求而形成的日志；
于Follower节点，日志项源自Leader节点的复制。
无论是Leader还是Follower，unstable log中的日志项最终都会被应用获取到并进行一系列处理（如写入WAL、存储至storage、
发送到其他Follower等），处理完成后，这些日志项可能就会变得不再有效，可以被回收。


### snapshot

[etcd-raft snapshot实现分析](https://zhuanlan.zhihu.com/p/29865583)

Follower节点被动接受Leader发送过来的snapshot后，需要将该snapshot应用到本身的状态机，其过程是：
Follower节点的raft内部状态机会将unstable log中的snapshot信息放在Ready结构中，应用通过Ready()
接口获取到snapshot信息，然后重放。


### leader transfer
[etcd raft如何实现leadership transfer](https://zhuanlan.zhihu.com/p/27895034)

大概原理是保证transferee(transfer的目标follower)拥有和原leader有一样新的日志，期间需要停写，
然后给transferee发送一个特殊的消息，让这个follower可以马上进行选主，而不用等到election timeout。
正常情况下，这个follower的term最大，当选，原来的leader变为follower。


http://int64.me/2017/Leader%20Transfer%20In%20TiKV.html


### Linearizable Read---ReadIndex Read

Linearizable Read, 通俗地讲就是读请求需要读到最新的已经commit的数据，不会读到老数据。
现实系统中，读请求通常会占很大比重，如果每次读请求都要走一次raft落盘，性能会大打折扣。

从raft协议可知，leader拥有最新的状态，如果读请求都走leader，那么leader可以直接返回结果给客户端。
然而，在出现网络分区和时钟快慢相差比较大的情况下，有可能会返回老的数据，即stale read。例如，leader
和其他followers之间出现网络分区，其他followers已经选出了新的leader，并且新的leader已经commit了一堆数据，
然而由于不同机器的时钟走的快慢不一，原来的leader可能并没有发觉自己的lease过期，仍然认为自己还是合法的
leader直接给客户端返回结果，从而导致了stale read。

Raft作者提出了一种叫做ReadIndex的方案：
当leader接收到读请求时，将当前commit index记录下来，记作read index，在返回结果给客户端之前，leader需要
先确定自己到底还是不是真的leader，确定的方法就是给其他所有peers发送一次心跳，如果收到了多数派的响应，
说明至少这个读请求到达这个节点时，这个节点仍然是leader，这时只需要等到commit index被apply到状态机后，
即可返回结果。

```go
func (n *node) ReadIndex(ctx context.Context, rctx []byte) error {
    return n.step(ctx, pb.Message{Type: pb.MsgReadIndex, Entries: []pb.Entry{{Data: rctx}}})
}
```
处理读请求时，应用的goroutine会调用这个函数，其中rctx参数相当于读请求id，全局保证唯一。step会往recvc中塞进一个MsgReadIndex消息，
而运行node入口函数run()的goroutine会从recvc中拿出这个message，并进行处理：
```
case m := <-n.recvc:
    // filter out response message from unknown From.
    if _, ok := r.prs[m.From]; ok || !IsResponseMsg(m.Type) {
        r.Step(m) // raft never returns an error
    }
```            
Step(m)最终会调用到raft结构体的step(m)，step是个函数指针，根据node的角色，运行stepLeader()/stepFollower()/stepCandidate()。

其中StepLeader的主要流程是
1.leader check自己是否在当前term commit过entry
2.leader会封装当前commit index和请求并保存到r.readOnly当中，然后给所有peers发心跳广播
3.当收到多数派响应时，再从r.readOnly中取出并以此新建一个ReadState结构体示例保存到readStates数组，
  而readStates数组会被包含在Ready结构体中
4.消费Ready的方法中会把readStates最后一个元素放入一个buffer为1的channal readStateC中
5.执行linearizableReadLoop()的goroutine从readStateC取出元素，首先判断请求的唯一id是否匹配，
然后等待apply index大于等于commit index，返回结果


https://zhuanlan.zhihu.com/p/27869566


### Linearizable Read---Lease Read

在 Raft 论文里面，提到了一种通过 clock + heartbeat 的 lease read 优化方法。 leader在发送 heartbeat 的时候，
会首先记录一个时间点 start，当系统大部分节点都回复了 heartbeat response，那么我们就可以认为 leader 的 lease 
有效期可以到 start + election timeout / clock drift bound 这个时间点。

为什么能够这么认为呢？主要是在于 Raft 的选举机制，因为 follower 会在至少 election timeout 的时间之后，才会重新发生选举，
所以下一个 leader 选出来的时间一定可以保证大于 start + election timeout / clock drift bound。

虽然采用 lease 的做法很高效，但仍然会面临风险问题，需要有一个前提条件，即各个服务器的 CPU clock 是准的，
即使有误差，也会在一个非常小的 bound 范围里面，如果各个服务器之间 clock 走的频率不一样，有些太快，有些太慢，
这套 lease 机制就可能出问题。

https://zhuanlan.zhihu.com/p/31118381
https://github.com/pingcap/blog-cn/blob/master/lease-read.md


[consensus-yaraft](https://blog.neverchanje.com/2017/08/03/consensus-yaraft/)
etcd raft 的c++实现


