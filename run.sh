run_auth_server() {
    cd ./build/auth
    ./auth start
}

run_data_server() {
    cd build/data
    ./data start
}

run_storage_server() {
    cd build/storage
    ./storage start
}

run_auth_db() {
    cd redis_conf
    redis-server auth.conf
}

run_data_db() {
    cd redis_conf
    redis-server data.conf
}

run_storage_db() {
    cd redis_conf
    redis-server storage.conf
}

case $1 in
    "server") 
        case $2 in
            "auth") run_auth_server;;
            "data") run_data_server;;
            "storage") run_storage_server;;
        esac
        ;;

    "db")
        case $2 in
            "auth") run_auth_db;;
            "data") run_data_db;;
            "storage") run_storage_db;;
        esac
        ;;
esac