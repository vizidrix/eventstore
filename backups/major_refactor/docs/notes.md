go test github.com/vizidrix/eventstore -bench .


http://grokbase.com/t/gg/golang-nuts/13499phn58/go-nuts-help-using-pprof
------------------------------------------------------------------------
go test -c
./package.test -test.run=XXX -test.bench=. -test.cpuprofile=cpu.out

go tool pprof package.test cpu.out

(pprof) top10
------------------------------------------------------------------------

./eventstore.text -test.run=1000 -test.bench=. -test.cpuprofile=cpu.out -test.memprofile=mem.out

./eventstore.test -test.bench=.20 -test.cpuprofile=cpu.out


 ASCII Subset

a-z (lower)		26
0-9 			10
					- 36 of 64 = 28 others

TODO: pick symbols for packed-ascii impl
-
_
@
\
/


Independent bits of command validation separated so it can be run at the boundaries

Make a memory map impl to test against also

* http://thinkbeforecoding.com/
Decide:
Command -> State -> Event list
ApplyStateChange:
State -> Event -> State


Constraints:
- Events are immutable
	- No need to accomodate updates
- Events are permanent
	- Need to accomodate physical deletion but only for repartitioning
	- Applications shouldn't have access only infrastructure layers
- Event data is opaque
	- No indexing beyond message header fields
	- No searching inside, view handlers are responsible for that

  Event Format

Folder per Domain
	-> Index file
		[ 8 byte Id + 4 byte Header offset + 4 byte Header CRC ]
		[ 8 byte header per event ]
	-> Data file per Aggregate Type partition (64Mb to start - 64Gb someday?)

Thoughts from LMDB:
Option for setting fixed key (required here) and also flag for fixed data size which could eliminate the LEN in header
xxx
Folder per Domain
	-> Folder per Aggregate Type
		-> File per Aggregate Id

V4:

-- Separated header info from data

[ Int16        | Int16       | Int32        ]
[ LEN 64k      | TYPE        | CRC          ]

[ byte[] ]
[ Data   ]

V3:
[ Int32 					| Int32 | byte[]	]
[ 3 bytes 		| 1 bytes 	|		|			]
[ LEN 4084 MAX  | EventType	| CRC	| DATA 		]

V2:
[ Int32 | Int32 | Int32                 			| byte[] ]
[ Int32 | Int32 | Byte	      | 3 Byte  			| byte[] ]
[ CRC	| TS	| EventType   | LEN 4084 MAX		| DATA	 ]

* Limit max len to 4096 - 12 to fit entry into 4kb page size



V1:
[ Int32 | Int32 | 2xInt64 							| Int32 | byte[] ]
[ CRC	| TS	| Domain 30 bytes + Type 2 bytes    | LEN 	| DATA	 ]
[ 		| 		| 20 packed char + 256 indexed ids	| 		|		 ]



Query by TS Range
Qeury by Namespace/EventType Set
Query by Index Range (position and size in set)?

id := es.NewKey()
registered := person.Registered("John", "Wayne", 987654321)
nameChanged := person.NameChanged("Jaughn", "Wayne")
profileUpdated := person.ProfileUpdated("Stuff has changed")
nameChangedReversed := person.NameChangeReversed("John", "Wayne")

es := eventstore.Connect()
domain := es.Domain("namespace")
aggregate := domain.Aggregate("person")
instance := aggregate.Instance(id)

err := instance.Put(registered)

- or -

err := instance.Put(
	registered,
	nameChanged,
	profileUpdated,
	nameChangeReversed)

events, err := instance.Get()


events, err := es.Get("namespace", "person", id)




New id's are seperated from existing id's
- New id's can be streamed in as a batch
- Existing id's can be split between quick append and realloc append
	- If there is room in the block then just append to existing data
	- If a data move is required then
		- Do the relocation(s) in batche(s)
		- Do the append into the newly available space


On any publish, if data size of batch will exceed targeted
block size (4k) for new id set then clamp the current gen 
and start writing to the next

http://www.cse.ohio-state.edu/~zhang/hpca11-submitted.pdf




// Write to ring -> 
//		calc CRC -> 
//			scan forward until
//			a) as far as possible with certain batch end reached
//				(must see next batch id to know current is finished)
//			b) memory buffer exceeds 4k buffer (roll back to last batch id)
//			- Write batch off to disk through data distribution logic
//			- Identify this batch as a "generation" to enable generational read/index concept
//		update index







