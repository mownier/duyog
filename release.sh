VERSION_FILE=./app/version.go
VERSION_CONTENT="
package app

// Version denotes the version of the program
const Version = \"%s\"
"

start() {
    rm -rf build
    ./build.sh

    if [ ! -d "release" ]; then
        mkdir release
    fi

    if [ ! -d "release/redis_conf" ]; then
        mkdir release/redis_conf
    fi

    cp -r ./build ./release
    cp ./redis_conf/auth.conf ./release/redis_conf/auth.conf
    cp ./redis_conf/data.conf ./release/redis_conf/data.conf
    cp ./redis_conf/storage.conf ./release/redis_conf/storage.conf
    cp ./db.sh ./release/duyog-db
    cp ./server.sh ./release/duyog-server
}

write_version() {
    echo "package app" > $VERSION_FILE
    echo "" >> $VERSION_FILE
    echo "// Version denotes the version of the program" >> $VERSION_FILE
    echo "const Version = \""$1"\"" >> $VERSION_FILE
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
    local version=$(extract_version)
    local major=$(extract_major $version)
    ((major++))
    version=$major.0
    write_version $version
}

update_minor() {
    local version=$(extract_version)
    local major=$(extract_major $version)
    local minor=$(extract_minor $version)
    ((minor++))
    version=$major.$minor
    write_version $version
}

update_bug_fix() {
    local version=$(extract_version)
    local major=$(extract_major $version)
    local minor=$(extract_minor $version)
    local bug_fix=$(extract_bug_fix $version)
    ((bug_fix++))
    version=$major.$minor.$bug_fix
    write_version $version
}

start

case $1 in
    "minor") update_minor;;
    "major") update_major;;
    "bug-fix") update_bug_fix;;
esac

echo v$(extract_version)
