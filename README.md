# Helm Chart Publisher
Helm Chart Publisher aims to help you build a nice CI/CD pipeline. It seats in front of a object storage service (such as AWS S3, OpenStack Swift) or a filesystem, sends your charts to it and also updates the index.

After receiving a PUT request with a repository and the chart, the publisher will upload the chart file to your storage, update the index and upload it too. Currently, it supports Amazon S3, OpenStack Swift, and Google Cloud Storage. Filesystem storage is planned.

## Configuration
The configuration is based on a YAML file. In order to publish your charts, you have to configure a `storage` and one or more `repos` (Helm repositories).

The Helm repository isolation can be done via bucket or a directory. The publisher will create an `index.yaml` for each repository you configure.

Each repo requires `name` and `bucket`. You can also specify a `directory`, if you do so, the charts are going to be stored in `bucket` under the specified path.

These are the configuration options for the helm publisher.

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
  gcs: {} # uses GCloud Application Default Credentials
  s3:
    bucket: charts-bucket-incubator
    region: us-west-2
  swift:
    username: SWIFT_USERNAME
    password: SWIFT_USERNAME
    tenant: some_tenant
    authUrl: https://some-auth-url:5000/v2.0
    endpointType: admin
    container: "kube-charts"
    insecureSkipVerify: false
```

## Running
You can either:
- Run as a binary
- Run as a Docker container

### Binary
Get the latest `helm-chart-publisher` for your platform on the [releases](https://github.com/luizbafilho/helm-chart-publisher/releases) page
```
curl -o /usr/local/bin/helm-chart-publisher -sSL https://github.com/luizbafilho/helm-chart-publisher/releases/download/<version>/helm-chart-publisher_<os>-<arch>
chmod +x /usr/local/bin/helm-chart-publisher
```

To run `helm-chart-publisher` you just have to execute the binary passing providing the configuration file.

```shell
$ PORT=8080 helm-chart-publisher --config /etc/helm-publisher/config.yaml
```

### Docker Container
To run the `helm-chart-publisher` from within a container run `make docker`. This will create the `helm-chart-publisher`
image.

Start the container using a command in the format:
```
docker run -p 8080:8080 -v /Users/$(whoami)/config:/etc/helm-chart-publisher/ -v /Users/$(whoami)/.aws:/root/.aws helm-chart-publisher
```

- The publisher runs on 8080 within the container and the ports must be exposed
- The config YAML must be mounted into the /etc/helm-chart-publisher directory
- .aws configuration should also be mounted in to prevent restarting container when credentials expire

### Inside minikube
To run helm-chart-publisher inside minikube:
- start minikube `minikube start`
- initialise helm `helm init`
- mount docker daemon from minikube `eval $(minikube docker-env)`
- run `make docker` to make the helm-chart-publisher image available to minikube
- edit values.yaml file to suit your needs
- if the storage option selected is AWS, run `minikube mount ~/.aws:/home/docker/.aws` this makes the credentials file available to
minikube to be mounted into the helm-chart-publisher container
- run `helm install --name publisher ./helm-chart-publisher` from the `helm` directory
- run `kubectl port-forward svc/publisher-helm-chart-publisher 8080:8080` helm-chart-publisher runs on port 8080,
forward the port to whatever localhost port you want
- test it is up and running using `curl http://localhost:8080/health`


## Usage 

You can publish a chart calling a simple `curl` command.

```
$ curl -i -X PUT -F repo=stable -F chart=@$HOME/charts/stable/mariadb-0.5.9.tgz http://localhost:8080/charts
```

This command will upload the chart file to an Amazon S3 bucket, updates the current `index.yaml` and upload it too.

The indexes are available via publisher under the `/:repo/index.yaml` path. For example to access the `stable` index.
```
$ curl -i http://localhost:8080/stable/index.yaml
```
But you can still access the `index.yaml` going directly to the storage. In this case:
```
$ curl -i https://s3-us-west-2.amazonaws.com/charts-bucket/index.yaml

apiVersion: v1
entries:
  mariadb:
  - created: 2017-03-07T17:36:04.782813678-03:00
    description: Fast, reliable, scalable, and easy to use open-source relational
      database system. MariaDB Server is intended for mission-critical, heavy-load
      production systems as well as for embedding into mass-deployed software.
    digest: d68c2852d7ac3e431cc65278d3ab74946b28479319a5707bc4405adf1dcd1393
    engine: gotpl
    home: https://mariadb.org
    icon: https://bitnami.com/assets/stacks/mariadb/img/mariadb-stack-220x234.png
    keywords:
    - mariadb
    - mysql
    - database
    - sql
    maintainers:
    - email: containers@bitnami.com
      name: Bitnami
    name: mariadb
    sources:
    - https://github.com/bitnami/bitnami-docker-mariadb
    urls:
    - https://s3-us-west-2.amazonaws.com/charts-bucket/mariadb-0.5.9.tgz
    version: 0.5.9
generated: 2017-03-07T17:34:47.965508312-03:00

```







## Roadmap
- [ ] Storages
  - [x] Openstack Swift
  - [x] Google Cloud Storage
  - [ ] Filesystem
- [ ] Tests
  - [ ] api
  - [ ] publisher
  - [ ] storages
  
  
## Notes
This project is at a very early stage, suggestions are, as always, very welcome in the form of PR's. If you feel the documentation is not clear or you have any questions, please open an issue for that.
