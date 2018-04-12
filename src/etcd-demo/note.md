

[etcd raft library设计原理和使用](http://www.cnblogs.com/foxmailed/p/7137431.html)

etcd提供的raft库只负责把消息保存在内存中，需要自行实现网络传输层以完成消息的发送和接收

raft提供了Ready结构体，其中封装了只读的entry和message, 可以持久化到存储或提交
或发送给其他peer，具体包括以下部分：
- pb.HardState：包含当前节点见过的最大的term，以及在这个term给谁投过票，已经当前节点知道的commit index
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
- 将Messages(上文提到的msgs)非阻塞的广播给其他peers
- 将CommittedEntries(已经commit还没有apply)应用到状态机。
- 如果发现CommittedEntries中有成员变更类型的entry，调用node的ApplyConfChange()进行通知
- 调用Node.Advance()告诉raft node，这批状态更新处理完了，可以处理下一批更新了。

应用通过raft.StartNode()来启动raft中的一个副本，函数内部通过启动一个goroutine运行
```go
func (n *node) run(r *raft)
```
来启动服务。


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

- propc: propc是一个没有buffer的channel，应用通过Propose方法写入的请求被封装成Message被push到propc中，
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


https://zhuanlan.zhihu.com/p/29180575


https://blog.neverchanje.com/2017/01/30/etcd_raft_core/


http://int64.me/2017/How%20to%20use%20Raft.html


https://medium.com/@daniel.chia/diving-into-etcd-raft-d48ce1cb6859


https://github.com/coreos/etcd/tree/master/contrib/raftexample



https://segmentfault.com/a/1190000008006649


