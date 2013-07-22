go test github.com/vizidrix/eventstore -bench .

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




  Event Format

Folder per Domain
	-> Folder per Aggregate Type
		-> File per Aggregate Id



[ Int32 | Int32 | 2xInt64 							| Int32 | byte[] ]
[ CRC	| TS	| Domain 30 bytes + Type 2 bytes    | LEN 	| DATA	 ]
[ 		| 		| 20 packed char + 256 indexed ids	| 		|		 ]




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
domain := es.Domain("WearShare")
aggregate := domain.Aggregate("Person")
instance := aggregate.Instance(id)

err := instance.Append(registered)

- or -

err := instance.Append(
	registered,
	nameChanged,
	profileUpdated,
	nameChangeReversed)

events, err := instance.ReadRaw()


events, err := es.ReadRaw("WearShare", "Person", id)





