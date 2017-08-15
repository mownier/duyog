build_auth_server() { 
    echo "Building auth server"

    if [ ! -d "build/auth" ]; then
        mkdir build/auth
    fi

    if go build ./auth/cli; then
        mv cli ./build/auth/auth
        cp ./auth/cli/config.json ./build/auth/config.json
    else
        exit
    fi
}

build_data_server() {
    echo "Building data server"
    
    if [ ! -d "build/data" ]; then
        mkdir build/data
    fi

    if go build ./data/cli; then
        mv cli ./build/data/data &&
        cp ./data/cli/config.json ./build/data/config.json
    else
        exit
    fi

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

    if go build ./storage/cli; then
        mv cli ./build/storage/storage &&
        cp ./storage/cli/config.json ./build/storage/config.json
    else
        exit
    fi
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

echo "Build finished"
