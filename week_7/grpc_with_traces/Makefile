LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-note-api
	make generate-other-note-api

generate-note-api:
	mkdir -p pkg/note_v1
	protoc --proto_path api/note_v1 \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/note_v1/note.proto

generate-other-note-api:
	mkdir -p pkg/other_note_v1
	protoc --proto_path api/other_note_v1 \
	--go_out=pkg/other_note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/other_note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/other_note_v1/other_note.proto

grpc-load-test:
	ghz \
		--proto api/note_v1/note.proto \
		--call note_v1.NoteV1.Get \
		--data '{"id": 1}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto api/note_v1/note.proto \
		--call note_v1.NoteV1.Get \
		--data '{"id": 0}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051