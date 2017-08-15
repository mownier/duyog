run_auth_server() {
    if ! ./build/auth/auth start -config ./build/auth/config.json; then
        exit
    fi
}

run_data_server() {
    if  ! ./build/data/data start -config ./build/data/config.json; then
        exit
    fi
}

run_storage_server() {
    if ./build/storage/storage start -config ./build/storage/config.json; then
        exit
    fi
}

run_all_servers() {
    trap exit_servers SIGINT
    run_auth_server >> ./build/auth.log 2<&1 &
    run_data_server >> ./build/data.log 2<&1 &
    run_storage_server >> ./build/storage.log 2<&1 &
    wait
}

exit_servers() {
    killall auth data storage
}

case $1 in
    "") 
        run_all_servers
        ;;

    "auth") 
        run_auth_server >> ./build/auth.log 2<&1
        ;;

    "data") 
        run_data_server >> ./build/data.log 2<&1
        ;;

    "storage")
        run_storage_server >> ./build/storage.log 2<&1
        ;;
esac
