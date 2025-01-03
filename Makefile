
.PHONY: all
all: api push

.PHONY: api
# generate api proto
api:
	buf dep update
	buf generate
	buf generate  --template buf.gen.example.yaml cmd/example  --debug 
.PHONY: push
push:
	(cd proto && buf push)
