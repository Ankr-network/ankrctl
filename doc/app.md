# Working with Ankr's application resource
`app` function is able to manage with all of your Ankr Network application resources. 

## Create a App:
To create an application, you need to choose the chart and namespace first, `chart list` and `chart detail` will help you locate the right chart, for namespace you can choose one from the `namespace list` output:
```
$ ankrctl app create testwp3 --chart-name=wordpress --chart-version=5.6.0 --chart-repo=stable --ns-id ns-1d8f3554-b678-4271-80b7-f72ab15e4f34
App app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5 create success.
```
or or create new one while you create the application:
```
$ ankrctl app create testwp4 --chart-name wordpress --chart-version 5.6.0 --chart-repo stable --ns-name testns4 --cpu-limit 300 --mem-limit 300 --storage-limit 5
App app-be31f9d3-858b-44d5-b314-68ea6a5a4582 create success.
```

## List all Apps:

```
$ ankrctl app list 

ID                                          Name       Chart Repo    Chart Name    Chart Version    App Version    Namespace    Cluster         Last Modify Date       Creation Date          Status         Event
app-6913b6e1-1c14-4096-98b7-d8a9d560b5a1    testwp1    stable        wordpress     5.6.0            5.1.0          wpns1        demo-cluster    10 May 19 14:55 PDT    10 May 19 14:55 PDT    app_running    launch_app_succeed
```

## List App Detail:

```
$ ankrctl app detail app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5

Application app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5 resource detail:
 LAST DEPLOYED: Sat May 11 19:47:46 2019
NAMESPACE: ns-1d8f3554-b678-4271-80b7-f72ab15e4f34
STATUS: DEPLOYED

RESOURCES:
==> v1/ConfigMap
NAME                                                    DATA  AGE
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-mariadb        1     6h38m
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-mariadb-tests  1     6h38m

==> v1/PersistentVolumeClaim
NAME                                                STATUS  VOLUME                                    CAPACITY  ACCESS MODES  STORAGECLASS  AGE
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-wordpress  Bound   pvc-a67d09f0-7425-11e9-a766-06dcd42c6c6c  10Gi      RWO           gp2           6h38m

==> v1/Secret
NAME                                                TYPE    DATA  AGE
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-mariadb    Opaque  2     6h38m
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-wordpress  Opaque  1     6h38m

==> v1/Service
NAME                                                TYPE          CLUSTER-IP      EXTERNAL-IP  PORT(S)                     AGE
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-mariadb    ClusterIP     100.64.189.200  <none>       3306/TCP                    6h38m
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-wordpress  LoadBalancer  100.71.128.34   <pending>    80:30493/TCP,443:31349/TCP  6h38m

==> v1beta1/Deployment
NAME                                                READY  UP-TO-DATE  AVAILABLE  AGE
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-wordpress  0/1    0           0          6h38m

==> v1beta1/StatefulSet
NAME                                              READY  AGE
app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-mariadb  0/1    6h38m


NOTES:
1. Get the WordPress URL:

  NOTE: It may take a few minutes for the LoadBalancer IP to be available.
        Watch the status with: 'kubectl get svc --namespace ns-1d8f3554-b678-4271-80b7-f72ab15e4f34 -w app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-wordpress'
  export SERVICE_IP=$(kubectl get svc --namespace ns-1d8f3554-b678-4271-80b7-f72ab15e4f34 app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-wordpress --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
  echo "WordPress URL: http://$SERVICE_IP/"
  echo "WordPress Admin URL: http://$SERVICE_IP/admin"

2. Login with the following credentials to see your blog

  echo Username: user
  echo Password: $(kubectl get secret --namespace ns-1d8f3554-b678-4271-80b7-f72ab15e4f34 app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5-wordpress -o jsonpath="{.data.wordpress-password}" | base64 --decode)
```

## Get App Overview:

```
$ ankrctl app overview
Cluster Count:		2
Namespace Count:	2
Network Count:		2
Total App Count:	4
Cluster Count:	2
Cpu Total:		800
Cpu Usage:	200
Mem Total:		800
Mem Usage:	266.66666
Storage Total:	15
Storage Usage:	7.5
```


## Update a App:

```
$ ankrctl app update app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5 --app-name testwp5 --update-version 5.7.1
App app-fb5f9f2c-40d4-4804-9e33-ee77beddeed5 update success.
```

## Cancel a App:

```
$ ankrctl app cancel app-be31f9d3-858b-44d5-b314-68ea6a5a4582
Warning: Are you sure you want to Cancel 1 app(s) (y/N) ? y
App app-be31f9d3-858b-44d5-b314-68ea6a5a4582 cancel success.
```
## Purge a App:

```
$ ankrctl app purge app-be31f9d3-858b-44d5-b314-68ea6a5a4582
Warning: Are you sure you want to Purge 1 app(s) (y/N) ? y
App app-be31f9d3-858b-44d5-b314-68ea6a5a4582 purge success.
```