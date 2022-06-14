# Mongo community operator

This is not offically supported helm charts. Testing this for local setup.

<https://github.com/mongodb/helm-charts/tree/main/charts/community-operator>

https://github.com/mongodb/helm-charts#unsupported-charts

Useful links:

<https://www.youtube.com/watch?v=VqeTT0NvRR4>
<https://www.youtube.com/watch?v=Pv70IcwipF0>

### Install steps

From the root of directory use this command to install the operator in kube.

```bash
make mongo
```

Step1: Install Operator for Mongo

Step2: Install resources for CRD

### Connect local mongo

Note: Install ingress for this connection.

```bash
mongo  --host localhost  --port 27017 -u my-user -p password
```

Enable Direct Connection option in Mongo Compass, so this will allow you login.

```
MongoDB shell version v4.2.0
connecting to: mongodb://localhost:27017/?compressors=disabled&gssapiServiceName=mongodb
Implicit session: session { "id" : UUID("2dcec1ca-3495-4a72-8f38-fe0580706013") }
MongoDB server version: 4.2.6
bpmn-mongodb:PRIMARY>
bye
```

Will document more hwo to connect this stuff here.

Use below connection string to connec the mongo running in kube. Use Primary DB to connect, secondary DB will not allow direct connection.

```bash
mongodb://my-user:password@mongo.example.com:27017/?directConnection=true&readPreference=primary&authMechanism=DEFAULT&authSource=admin&retryWrites=false
```
