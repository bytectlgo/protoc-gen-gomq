
.PHONY: api
# generate api proto
api:
	buf generate
	buf generate  --template buf.gen.example.yaml example  --debug 
	 