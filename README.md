# Helm Chart Publisher
---
Helm Chart Publisher aims to help you build a nice CI/CD pipeline. It seats in front of a storage and send your charts to it and also updates the index.

After receive a PUT request with a repository and the chart, the publisher will upload the chart file to your storage, updates the index and upload it too. Currently supports only Amazon S3, but OpenStack Swift, Google Cloud Storage and Filesystem are planned.

## Configuration
The configuration is based on a YAML file. In order to publish your charts, you have to configure a `storage` and one or more `repos` (Helm repositories).

The Helm repository isolation can be done via bucket or directory. The publisher will create an `index.yaml` for each repository you configure.

Each repo requires `name` and `bucket`. You can also specify a `directory`, if you do so, the charts are going to be stored in `bucket` under the specified path.

These are configuration options for the helm publisher.

```
repos:
  - name: stable
    bucket: charts-bucket
  - name: incubator
    bucket: charts-bucket-incubator
  - name: test
    bucket: test-bucket
    directory: test

storage:
  s3:
    accessKey: AMAZON_ACCESS_KEY
    secretKey: AMAZON_SECRET_KEY
    region: us-west-2
```

## Usage

To run `helm-chart-publisher` you just have to execute the binary passing the configuration file.

```shell
$ PORT=8080 helm-chart-publisher --config /etc/helm-publisher/config.yaml
```

You can publish a chart very easily by call a simple `curl`

```
curl -i -X PUT -F repo=stable -F chart=@$HOME/charts/stable/mariadb-0.5.9.tgz http://localhost:8080/charts
```

This command will upload the chart file to a Amazon S3 bucket, updates the current `index.yaml` and upload it too.

The indexes are available via publisher under `/:repo/index.yaml` path. For example to access the `stable` index.
```
$ curl -i http://localhost:8080/stable/index.yaml
```
