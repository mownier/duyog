AUTH_SERVER=./build/auth/auth
AUTH_SERVER_PID=./build/auth/auth.pid
AUTH_SERVER_LOG=./build/auth/auth.log
AUTH_SERVER_CNF=./build/auth/config.json

DATA_SERVER=./build/data/data
DATA_SERVER_PID=./build/data/data.pid
DATA_SERVER_LOG=./build/data/data.log
DATA_SERVER_CNF=./build/data/config.json

STORAGE_SERVER=./build/storage/storage
STORAGE_SERVER_PID=./build/storage/storage.pid
STORAGE_SERVER_LOG=./build/storage/storage.log
STORAGE_SERVER_CNF=./build/storage/config.json

start_auth_server() {
    if [ -f $AUTH_SERVER_PID ]; then
        echo "already running auth server"
    else
        $AUTH_SERVER start -config $AUTH_SERVER_CNF >> $AUTH_SERVER_LOG 2<&1 & auth_pid=$!
        echo "auth server started"
        echo $auth_pid > $AUTH_SERVER_PID
    fi
}

start_data_server() {
    if [ -f $DATA_SERVER_PID ]; then
        echo "already running data server"
    else
        $DATA_SERVER start -config $DATA_SERVER_CNF >> $DATA_SERVER_LOG 2<&1 & data_pid=$!
        echo "data server started"
        echo $data_pid > $DATA_SERVER_PID
    fi
}

start_storage_server() {
    if [ -f $STORAGE_SERVER_PID ]; then
        echo "already running storage server"
    else
        $STORAGE_SERVER start -config $STORAGE_SERVER_CNF >> $STORAGE_SERVER_LOG 2<&1 & storage_pid=$!
        echo "storage server started"
        echo $storage_pid > $STORAGE_SERVER_PID
    fi
}

stop_auth_server() {
    if [ ! -f $AUTH_SERVER_PID ]; then
        echo "no auth server process running"
    else
        pkill -F $AUTH_SERVER_PID
        rm $AUTH_SERVER_PID
        echo "auth server stopped"
    fi
}

stop_data_server() {
    if [ ! -f $DATA_SERVER_PID ]; then
        echo "no data server process running"
    else
        pkill -F $DATA_SERVER_PID
        rm $DATA_SERVER_PID
        echo "data server stopped"
    fi
}

stop_storage_server() {
    if [ ! -f $STORAGE_SERVER_PID ]; then
        echo "no storage server process running"
    else
        pkill -F $STORAGE_SERVER_PID
        rm $STORAGE_SERVER_PID
        echo "storage server stopped"
    fi
}

start_all_servers() {
    start_auth_server &&
    start_data_server &&
    start_storage_server
}

stop_all_servers() {
    stop_storage_server &&
    stop_data_server &&
    stop_auth_server
}

list_servers() {
    local auth_s=$(is_server_running $AUTH_SERVER_PID)
    local data_s=$(is_server_running $DATA_SERVER_PID)
    local stor_s=$(is_server_running $STORAGE_SERVER_PID)
    echo "+---------+-----------+"
    echo "| server  | status    |"
    echo "+---------+-----------+"
    echo "| auth    | $auth_s   |"
    echo "| data    | $data_s   |"
    echo "| storage | $stor_s   |"
    echo "+---------+-----------+"
}

is_server_running() {
    if [ -f $1 ]; then
        echo "started"
    else
        echo "stopped"
    fi
}

case $1 in
    "start") 
        case $2 in
            "") start_all_servers;;
            "auth") start_auth_server;;
            "data") start_data_server;;
            "storage") start_storage_server;;
        esac;;
    
    "stop")
        case $2 in
            "") stop_all_servers;;
            "auth") stop_auth_server;;
            "data") stop_data_server;;
            "storage") stop_storage_server;;
        esac;;

    "list")
        list_servers
        ;;
esac
