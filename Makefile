
.PHONY: protobuf
# generate api proto
protobuf:
	protoc --proto_path=./ \
			--proto_path=./third_party \
			--go_out=paths=source_relative:.\
			protobuf/mq.proto

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./ \
			--proto_path=./protobuf \
			--proto_path=./third_party \
			--go_out=paths=source_relative:.\
			--gomq_out=paths=source_relative:.\
			example/example.proto

	 