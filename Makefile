generate-go-messages:	
					protoc --go_out=server/pb --go_opt=paths=source_relative \
   					--go-grpc_out=server/pb --go-grpc_opt=paths=source_relative \
    				messages/messages.proto