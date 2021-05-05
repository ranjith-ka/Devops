# ProtoBuffers

### Compiler setup

To install `protoc` use make command utilities.

```bash
make install
```
To manage the version which I choose, if not use

`go get -u github.com/golang/protobuf/proto`

The below command used to generate the go file from protoc

```bash
protoc -I proto/ --go_out=. proto/03_city_message.proto
```
