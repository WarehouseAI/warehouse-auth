proto_dir="./proto"
proto_out="./internal/warehousepb"

if [ ! -d "$proto_dir" ]; then
  echo "Directory '$proto_dir' does not exist"
  exit 1
fi

rm -fr $proto_out
mkdir $proto_out
for file in "$proto_dir"/*.proto; do
  if [ -f "$file" ]; then
    echo "Compiling $file..."
    protoc --proto_path="$proto_dir" --go_out=./internal --go-grpc_out=./internal "$file"
  fi
done

echo "Completed"