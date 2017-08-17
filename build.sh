AUTH_BUILD_FILE=./auth/app/build.go
AUTH_VERSION_FILE=./auth/app/version.go

DATA_BUILD_FILE=./data/app/build.go
DATA_VERSION_FILE=./data/app/version.go

STORAGE_BUILD_FILE=./storage/app/build.go
STORAGE_VERSION_FILE=./storage/app/version.go

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

print_usage() {
    echo ""
    echo "USAGE:"
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
    echo "package app" > $2
    echo "" >> $2
    echo "// Build denotes the build of the program" >> $2
    echo "const Build = \""$1"\"" >> $2
}

update_build() {
    local version=$(extract_version $1)
    local major=$(extract_major $version)
    local new=${major}D$(date +"%y")U$(date +"%V")Y$(date +"%u")O$(date +"%H")
    local build=$(extract_build $2)
    local b=`echo $build | cut -d \G -f 2`
    ((b++))
    new=${new}G${b}
    write_build $new $2
}

extract_build() {
    local content=$(cat $1)
    local b=${content:68}
    b=${b/\"/""}
    b=${b/\"/""}
    echo $b
}

extract_version() {
    local content=$(cat $1)
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
else
    if [ "$1" = "" ]; then
        exit 0
    fi
    
    if [ ! -d "build" ]; then
        mkdir build
    fi

    case $1 in
        "auth") 
            update_build $AUTH_VERSION_FILE $AUTH_BUILD_FILE && build_auth_server
            echo auth server new build: $(extract_build $AUTH_BUILD_FILE);;

        "data") 
            update_build $DATA_VERSION_FILE $DATA_BUILD_FILE && build_data_server
            echo data server new build: $(extract_build $DATA_BUILD_FILE);;

        "storage") 
            update_build $STORAGE_VERSION_FILE $STORAGE_BUILD_FILE && build_storage_server
            echo storage server new build: $(extract_build $STORAGE_BUILD_FILE);;
    esac

    echo "build finished"
fi

