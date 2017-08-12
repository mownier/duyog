rm -rf build
mkdir build
mkdir build/auth
mkdir build/data
mkdir build/storage
mkdir build/storage/file
mkdir build/storage/file/song
mkdir build/storage/file/user
mkdir build/storage/file/album
mkdir build/storage/file/artist
mkdir build/storage/file/playlist

echo "Building auth server"
cd auth/cli/
go build
mv cli ../../build/auth/auth
cp config.json ../../build/auth/config.json

echo "Building data server"
cd ../../data/cli
go build
mv cli ../../build/data/data
cp config.json ../../build/data/config.json

echo "Building storage server"
cd ../../storage/cli
go build
mv cli ../../build/storage/storage
cp config.json ../../build/storage/config.json

echo "Finished"
