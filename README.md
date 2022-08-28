# Go Http Over GRPC

### How to Install Protocol Buffers on Windows?
- At first, browser https://github.com/protocolbuffers/protobuf/releases
- download protoc-xxxxx-win32.zip or protoc-xxxxx-win64.zip
- extract and save in C:\
- set file C:\protoc-xxxxx-win64\bin in environment variabel ‘System variables’ in variabel path

### How to Install Buf on Windows?
- Browse https://docs.buf.build/installation
- Download Buf for Windows and rename file become buf.exe
- route file buf.exe to environment variabel ‘System variables’ in variabel path
- Create file buf.gen.yml
- execute command `buf mod update`
- Always run buf mod update after adding a dependency to your buf.yaml

### Generate File Proto (Old Way)
- to folder /proto
- exec ./get-googleapi.sh (Only 1 exec for create library)
- create file xxxxx.proto in v1
- exec cmd `protoc --proto_path=. v1/*.proto --go_out=plugins=grpc:./ --grpc-gateway_out=:./`

### Generate File Proto (New Way)
- exec command buf generate

### Running apps
- go run .



### References
- https://cloud.google.com/endpoints/docs/grpc/transcoding#configuring_transcoding_in_yaml
- https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#google.api.HttpRule
- https://grpc-ecosystem.github.io/grpc-gateway
- https://github.com/grpc-ecosystem/grpc-gateway
- https://github.com/johanbrandhorst/grpc-gateway-boilerplate
- https://github.com/bufbuild/buf
- https://www.freecodecamp.org/news/a-quick-introduction-to-clean-architecture-990c014448d2/