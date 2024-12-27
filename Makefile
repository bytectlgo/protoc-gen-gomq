
.PHONY: api
# generate api proto
api:
	buf dep update
	buf generate
	buf generate  --template buf.gen.example.yaml cmd/example  --debug 
	 