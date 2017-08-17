BUILD_FILE=./app/build.go
VERSION_FILE=./app/version.go

build_auth_server() { 
    echo "building auth server"

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
    echo "building data server"
    
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
    echo "building storage server"

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

print_usage() {
    echo ""
    echo "USAGE:"
    echo ""
    echo "  note: '$ ./build.sh' will build all the servers and update the build file"
    echo ""
    echo "  help"
    echo "    - Prints the usage info"
    echo ""
    echo "  auth"
    echo "    - Builds the auth server"
    echo ""
    echo "  data"
    echo "    - Builds the data server"
    echo ""
    echo "  storage"
    echo "    - Builds the storage server"
    echo ""
}

write_build() {
    echo "package app" > $BUILD_FILE
    echo "" >> $BUILD_FILE
    echo "// Build denotes the build of the program" >> $BUILD_FILE
    echo "const Build = \""$1"\"" >> $BUILD_FILE
}

update_build() {
    local version=$(extract_version)
    local major=$(extract_major $version)
    local new=${major}D$(date +"%y")U$(date +"%V")Y$(date +"%u")O$(date +"%H")
    local build=$(extract_build)
    local b=`echo $build | cut -d \G -f 2`
    ((b++))
    write_build ${new}G${b}
}

extract_build() {
    local content=$(cat $BUILD_FILE)
    local b=${content:68}
    b=${b/\"/""}
    b=${b/\"/""}
    echo $b
}

extract_version() {
    local content=$(cat $VERSION_FILE)
    local v=${content:74}
    v=${v/\"/""}
    v=${v/\"/""}
    echo $v
}

extract_major() {
    local v=`echo $1 | cut -d \. -f 1`
    if [ ${#v} = 0 ]; then
        echo 0
    else
        echo $v
    fi
}

if [ "$1" = "help" ]; then
    print_usage
    update_build
else
    if [ ! -d "build" ]; then
        mkdir build
    fi

    case $1 in
        "") build_all_servers && update_build;;
        "auth") build_auth_server;;
        "data") build_data_server;;
        "storage") build_storage_server;;
    esac

    echo "build finished"
    echo build: $(extract_build)
fi

