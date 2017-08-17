AUTH_VERSION_FILE=./auth/app/version.go
DATA_VERSION_FILE=./data/app/version.go
STORAGE_VERSION_FILE=./storage/app/version.go

start_build() {
    rm -rf build/$1
    ./build.sh $1
}

write_version() {
    echo "package app" > $2
    echo "" >> $2
    echo "// Version denotes the version of the program" >> $2
    echo "const Version = \""$1"\"" >> $2
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

extract_minor() {
    local v=`echo $1 | cut -d \. -f 2`
    if [ ${#v} = 0 ]; then
        echo 0
    else
        echo $v
    fi
}

extract_bug_fix() {
    local v=`echo $1 | cut -d \. -f 3`
    if [ ${#v} = 0 ]; then
        echo 0
    else
        echo $v
    fi
}

update_major() {
    local version=$(extract_version $1)
    local major=$(extract_major $version)
    ((major++))
    version=$major.0
    write_version $version $1
}

update_minor() {
    local version=$(extract_version $1)
    local major=$(extract_major $version)
    local minor=$(extract_minor $version)
    ((minor++))
    version=$major.$minor
    write_version $version $1
}

update_bug_fix() {
    local version=$(extract_version $1)
    local major=$(extract_major $version)
    local minor=$(extract_minor $version)
    local bug_fix=$(extract_bug_fix $version)
    ((bug_fix++))
    version=$major.$minor.$bug_fix
    write_version $version $1
}

print_usage() {
    echo ""
    echo "USAGE:"
    echo ""
    echo "  help"
    echo "    - Prints the usage info"
    echo ""
    echo "  [auth|data|storage] [major|minor|bug-fix]"
    echo ""
    echo "  major"
    echo "    - Increment the major version of a server"
    echo "    - If current version is 2.1.1, it will update to 3.0"
    echo ""
    echo "  minor"
    echo "    - Increment the minor version of a server"
    echo "    - If current version is 2.0.1, it will update to 2.1"
    echo ""
    echo "  bug-fix"
    echo "    - Increment the bug-fix version of a server"
    echo "    - If current version is 2.3.4, it will update to 2.3.5"
    echo ""
}

if [ "$1" = "help" ]; then
    print_usage
else
    if [ "$2" = "" ]; then
        exit 0
    fi

    if [ ! -d "release" ]; then
        mkdir release
    fi

    if [ ! -d "release/build" ]; then
        mkdir release/build
    fi

    if [ ! -d "release/redis_conf" ]; then
        mkdir release/redis_conf
    fi

    case $1 in
        "auth")
            case $2 in
                "minor") update_minor $AUTH_VERSION_FILE;;
                "major") update_major $AUTH_VERSION_FILE;;
                "bug-fix") update_bug_fix $AUTH_VERSION_FILE;;
            esac

            start_build "auth"

            if [ ! -d "release/build/auth" ]; then
                mkdir release/build/auth
            fi

            cp -rf ./build/auth ./release/build/
            cp ./redis_conf/auth.conf ./release/redis_conf/auth.conf
            echo auth server new version: $(extract_version $AUTH_VERSION_FILE);;
        
        "data")
            case $2 in
                "minor") update_minor $DATA_VERSION_FILE;;
                "major") update_major $DATA_VERSION_FILE;;
                "bug-fix") update_bug_fix $DATA_VERSION_FILE;
            esac

            start_build "data"

            if [ ! -d "release/build/data" ]; then
                mkdir release/build/data
            fi

            cp -rf ./build/data ./release/build/
            cp ./redis_conf/data.conf ./release/redis_conf/data.conf
            echo data server new version: $(extract_version $DATA_VERSION_FILE);;

        "storage")
            case $2 in
                "minor") update_minor $STORAGE_VERSION_FILE;;
                "major") update_major $STORAGE_VERSION_FILE;;
                "bug-fix") update_bug_fix $STORAGE_VERSION_FILE;;
            esac

            start_build "storage"

            if [ ! -d "release/build/storage" ]; then
                mkdir release/build/storage
            fi

            cp -rf ./build/storage ./release/build/
            cp ./redis_conf/storage.conf ./release/redis_conf/storage.conf
            echo storage server new version: $(extract_version $STORAGE_VERSION_FILE);;
    esac

    if [ ! -f "release/duyog-db" ]; then
        cp db.sh ./release/duyog-db
    fi

    if [ ! -f "release/duyog-server" ]; then
        cp server.sh ./release/duyog-server
    fi
fi
