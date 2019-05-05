# Working with Ankr's chart resource
`chart` function is able to manage with all of your Ankr Network chart resources. 

## List all Charts:

```
$ ankrctl chart list --list-repo stable
Repo      Name         Latest Version    Latest App Version    Description
stable    wordpress    5.7.1             5.1.1                 Web publishing platform for building blogs and websites.
```

## Upload a Chart:
Upload a new chart to user catalog:
```
$ ankrctl upload wordpress --upload-file ../dccn-appmgr/examples/test/wordpress-5.7.1.tgz --upload-version=5.7.1
Chart wordpress upload success.
```

## Delete a Chart:

```
$ ankrctl chart delete wordpress --delete-version=5.7.1
Warning: Are you sure you want to Delete chart wordpress version 5.7.1 (y/N) ? y
Chart wordpress version 5.7.1 delete success.
```

## List Chart detail:

```
$ ankrctl chart detail wordpress --detail-repo stable --show-version 5.6.0
Repo: stable	Chart: wordpress
Version		App Version
5.7.1		5.1.1
5.6.0		5.1.0

++++++++++ Chart versions 5.6.0 readme.md ++++++++++
<output of readme.md>
++++++++++ Chart versions 5.6.0 values.yaml ++++++++++
<output of values.yaml>
```
## Save as a new chart version with updated values.yaml:

```
$ ankrctl chart saveas wordpress-5.7.2 --saveas-version 5.7.2 --source-name wordpress --source-repo stable --source-version 5.7.1 --values-yaml ./values.yaml
Chart wordpress-5.7.2 save success.
```
