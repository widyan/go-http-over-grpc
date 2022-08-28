# Go Http Over GRPC

### How to Install Protocol Buffers on Windows?
- At first, browser https://github.com/protocolbuffers/protobuf/releases
- Download protoc-xxxxx-win32.zip or protoc-xxxxx-win64.zip
- Extract and save in C:\
- Set file C:\protoc-xxxxx-win64\bin in environment variabel ‘System variables’ in variabel path

### How to Install Buf on Windows?
- Browse https://docs.buf.build/installation
- Download Buf for Windows and rename file become buf.exe
- Route file buf.exe to environment variabel ‘System variables’ in variabel path
- Create file buf.gen.yml
- Execute command `buf mod update`
- Always run buf mod update after adding a dependency to your buf.yaml

### Generate File Proto (Old Way)
- To folder /proto
- Exec ./get-googleapi.sh (Only 1 exec for create library)
- Create file xxxxx.proto
- Exec cmd `protoc --proto_path=. v1/*.proto --go_out=plugins=grpc:./ --grpc-gateway_out=:./`

### Generate File Proto and generate OpenAPI in folder third party (New Way)
- Exec command `buf generate`

### Running apps
- Add .env (reference from .env.example)
- Exec command `go run .`
- If success, run http://localhost:7000 swagger URI from browser


### References
- https://cloud.google.com/endpoints/docs/grpc/transcoding#configuring_transcoding_in_yaml
- https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#google.api.HttpRule
- https://grpc-ecosystem.github.io/grpc-gateway
- https://github.com/grpc-ecosystem/grpc-gateway
- https://github.com/johanbrandhorst/grpc-gateway-boilerplate
- https://github.com/bufbuild/buf
- https://www.freecodecamp.org/news/a-quick-introduction-to-clean-architecture-990c014448d2/