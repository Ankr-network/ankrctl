# Working with Ankr's namespace resource
`namespace` function is able to manage with all of your Ankr Network namespace resources. 

## List all Namespaces:

```
$ ankrctl namespace list
ID                                         Name       CPU Limit      Memory Limit    Storage Limit    Cluster ID                                     Cluster Name    Status            Event
ns-1d8f3554-b678-4271-80b7-f72ab15e4f34    wpns1      0.5 vCPU(s)    0.5 GB          10 GB            daemon-44f9477a-c39c-4107-a880-edb98e566e51    demo-cluster    NS_RUNNING        LAUNCH_NS_SUCCEED
```

## Create a Namespace:
```
$ ankrctl namespace create testns01 --cpu-limit 200 --mem-limit 200 --storage-limit 3
Namespace ns-89bb94be-7e7f-4f6e-8f5e-e3103c5ed68c create success.
```

## Delete a Namespace:

```
$ ankrctl namespace delete ns-5f9e00af-025d-482f-812d-01c488cd70b9
Warning: Are you sure you want to Cancel 1 namespace(s) (y/N) ? y
Namespace ns-5f9e00af-025d-482f-812d-01c488cd70b9 delete success.
```

## Update a Namespace:

```
$ ankrctl update ns-89bb94be-7e7f-4f6e-8f5e-e3103c5ed68c --cpu-limit 300 --mem-limit 300 --storage-limit 4
Namespace ns-89bb94be-7e7f-4f6e-8f5e-e3103c5ed68c update success.
```
