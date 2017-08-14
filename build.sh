build_auth_server() { 
    echo "Building auth server"

    if [ ! -d "build/auth" ]; then
        mkdir build/auth
    fi

    go build ./auth/cli
    mv cli ./build/auth/auth
    cp ./auth/cli/config.json ./build/auth/config.json
}

build_data_server() {
    echo "Building data server"
    
    if [ ! -d "build/data" ]; then
        mkdir build/data
    fi

    go build ./data/cli
    mv cli ./build/data/data
    cp ./data/cli/config.json ./build/data/config.json
}

build_storage_server() {
    echo "Building storage server"

    if [ ! -d "build/storage" ]; then
        mkdir build/storage
    fi

    if [ ! -d "build/storage/file" ]; then
        mkdir build/storage/file
    fi

    if [ ! -d "build/storage/file/song" ]; then
        mkdir build/storage/file/song
    fi

    if [ ! -d "build/storage/file/user" ]; then
        mkdir build/storage/file/user
    fi

    if [ ! -d "build/storage/file/album" ]; then
        mkdir build/storage/file/album
    fi

    if [ ! -d "build/storage/file/artist" ]; then
        mkdir build/storage/file/artist
    fi

    if [ ! -d "build/storage/file/playlist" ]; then
        mkdir build/storage/file/playlist
    fi

    go build ./storage/cli
    mv cli ./build/storage/storage
    cp ./storage/cli/config.json ./build/storage/config.json
}

build_all_servers() {
    build_auth_server
    build_data_server
    build_storage_server
}

if [ ! -d "build" ]; then
    mkdir build
fi


case $1 in
    "") build_all_servers;;
    "auth") build_auth_server;;
    "data") build_data_server;;
    "storage") build_storage_server;;
esac

echo "Finished"
