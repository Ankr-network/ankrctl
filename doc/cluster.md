# Working with Ankr's cluster resource
`cluster` function is able to manage with all of your Ankr Network cluster resources. 

##  List all Cluster:

```
$ ankrctl cluster list
ID                                             Name            CPU        Memory    Storage     Latitude     Longitude      Status       WalletAddress
daemon-44f9477a-c39c-4107-a880-edb98e566e51    demo-cluster    6CPU(s)    5.47GB    121.32GB    37.774900    -122.419400    available
```

##  List Cluster Network Info:

```
$ ankrctl cluster network
User Count:		299
Host Count:		137
Namespace Count:	450
Container Count:	1342
Traffic:	1
```