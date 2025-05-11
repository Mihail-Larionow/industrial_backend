$protoc = "protoc"
$protoFile = "api/proto/calculator.proto"
$goOut = "--go_out=."
$goOpt = "--go_opt=paths=source_relative"
$goGrpcOut = "--go-grpc_out=."
$goGrpcOpt = "--go-grpc_opt=paths=source_relative"

& $protoc $goOut $goOpt $goGrpcOut $goGrpcOpt $protoFile 