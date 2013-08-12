

Very interesting CQRS Blogger:
http://thinkbeforecoding.com/

BitCask:
https://github.com/basho/bitcask/

NoSQL in Golang:
https://github.com/HouzuoGuo/tiedot

Queue comparison (Redis and ZMQ) in Golang:
(Wonder what LMDB could do for redis like space?)
https://github.com/stephenmcd/two-queues

-> ZMQ in Golang: https://github.com/alecthomas/gozmq

Some fun MIT Open Courseware videos:
http://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-172-performance-engineering-of-software-systems-fall-2010/video-lectures/


Howard Chu - Author of LMDB and related Source:
https://gitorious.org/~hyc
- OMG fast database... wondering if event store could be a constrained version of this?
- How important is it to keep events in contiguous blocks on disk?
- Could append only scheme work without compaction or would that degrade over time?
- Use actual LMDB store to hold aggregate snapshots, indexes and views?

LMDB library for Golang:
https://github.com/szferi/gomdb

HyperDex:
Affiliated with Howard Chu?
Might already solve the challenges for distributed coordination in leu of trying to build it on top of LMDB and Raft...?
http://hyperdex.org/papers/


Raft Distributed Election Protocol:
http://www.youtube.com/watch?v=YbZ3zDzDnrw
- Could something like this be used to coordinate distributed Read among peers?

i.e. Quorum amongst peers to determine latest 'servable' event version
- Natural partitions at Domain, Kind and Id could make manageable sized 'elected leader of partition' blocks?
- Dead nodes would be easy to handle, just sync on re-connect
- Network fragmentation could rely on lack of Quorum in heartbeats to quickly stop sending stale data, only majority partition could continue?
- Even stale data would be valid, just out of date due to event sourcing




Chainable Replication (Similar to HyperDex?):
http://muratbuffalo.blogspot.com/2011/02/chain-replication-for-supporting-high.html


Haystack (Facebook) Image Store:
http://muratbuffalo.blogspot.com/2010/12/finding-needle-in-haystack-facebooks.html

Ceph:
http://muratbuffalo.blogspot.com/2011/03/ceph-scalable-high-performance.html