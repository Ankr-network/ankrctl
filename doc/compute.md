# Working with Ankr's computing resource
`compute` function is able to interact with all of your Ankr Network distributed cloud computing network resources. 
## Task Commands
### List all Tasks:

```
$ ankrctl compute task list 

Task Id                                 Task Name    Type          Image         Last Modify Date       Creation Date          Replica    Data Center         Status
b3debe9d-dc82-4560-aaa6-7f4b445f7566    test001      deployment    nginx:1.12    05 Mar 19 06:57 UTC    05 Mar 19 06:57 UTC    2          datacenter_tokyo    running
```
### Create a Task:
```
$ ankrctl compute task create test001 --image=nginx:1.12 --replica 2

Task 9f0386f8-51b9-43d0-a624-adf39a512b87 Create Success.
```
### Cancel a Task:

```
$ ankrctl compute task cancel 9f0386f8-51b9-43d0-a624-adf39a512b87

Warning: Are you sure you want to Cancel 1 task(s) (y/N) ? y
Task 9f0386f8-51b9-43d0-a624-adf39a512b87 Cancel Success.
```
### Purge a Task:

```
$ ankrctl compute task purge 9f0386f8-51b9-43d0-a624-adf39a512b87

Warning: Are you sure you want to Purge 1 task(s) (y/N) ? y
Task 9f0386f8-51b9-43d0-a624-adf39a512b87 Purge Success.
```
### Update a Task:

```
$ ankrctl compute task update b3debe9d-dc82-4560-aaa6-7f4b445f7566 --image nginx:1.13 --replica 1

Task b3debe9d-dc82-4560-aaa6-7f4b445f7566 Update Success.
```
## Working with Datacenter
###  List all Datacenter:

```
$ ankrctl compute dc list

Id                                      Name                    CPU        RAM        HDD         Latitude    Longitude    Status       WalletAddress
45828aa9-9b45-4f40-9c07-47182ad512a5    datacenter_tokyo        6CPU(s)    9.24GB     161.82GB    35.6850     139.7510     available
5ef4f226-a655-46a9-92aa-d6c1101045a9    datacenter_singapore    8CPU(s)    4.37GB     215.76GB    1.2931      103.8560     available
0a50e2d5-671b-429a-a9ab-4578caef6345    datacenter_seoul        4CPU(s)    3.65GB     107.88GB    37.5985     126.9780     available
eee70d94-688c-4984-b915-88d3855ffebc    datacenter_portland     6CPU(s)    11.12GB    323.80GB    43.7442     -120.3890    available
```