REDIS=redis-server

AUTH_REDIS_CNF=./auth.conf
AUTH_REDIS_DUYOG_PID=./auth_duyog.pid

DATA_REDIS_CNF=./data.conf
DATA_REDIS_DUYOG_PID=./data_duyog.pid

STORAGE_REDIS_CNF=./storage.conf
STORAGE_REDIS_DUYOG_PID=./storage_duyog.pid

start_auth_redis() {
    if [ -f $AUTH_REDIS_DUYOG_PID ]; then
        echo "already running redis db for auth server"
    else
        $REDIS $AUTH_REDIS_CNF & d_auth_pid=$!
        echo "started redis for auth server"
        echo $d_auth_pid > $AUTH_REDIS_DUYOG_PID
    fi
}

start_data_redis() {
    if [ -f $DATA_REDIS_DUYOG_PID ]; then
        echo "already running redis for data server"
    else
        $REDIS $DATA_REDIS_CNF & d_data_pid=$!
        echo "started redis for data server"
        echo $d_data_pid > $DATA_REDIS_DUYOG_PID
    fi
}

start_storage_redis() {
    if [ -f $STORAGE_REDIS_DUYOG_PID ]; then
        echo "already running redis for storage server"
    else
        $REDIS $STORAGE_REDIS_CNF & d_storage_pid=$!
        echo "started redis for storage server"
        echo $d_storage_pid > $STORAGE_REDIS_DUYOG_PID
    fi
}

stop_auth_redis() {
    if [ ! -f $AUTH_REDIS_DUYOG_PID ]; then
        echo "no auth redis process running"
    else
        pkill -F $AUTH_REDIS_DUYOG_PID
        rm $AUTH_REDIS_DUYOG_PID
        echo "auth redis stopped"
    fi
}

stop_data_redis() {
    if [ ! -f $DATA_REDIS_DUYOG_PID ]; then
        echo "no data redis process running"
    else
        pkill -F $DATA_REDIS_DUYOG_PID
        rm $DATA_REDIS_DUYOG_PID
        echo "data redis stopped"
    fi
}

stop_storage_redis() {
    if [ ! -f $STORAGE_REDIS_DUYOG_PID ]; then
        echo "no storage redis process running"
    else
        pkill -F $STORAGE_REDIS_DUYOG_PID
        rm $STORAGE_REDIS_DUYOG_PID
        echo "storage redis stopped"
    fi
}

start_all_redis() {
    start_auth_redis &&
    start_data_redis &&
    start_storage_redis
}

stop_all_redis() {
    stop_auth_redis &&
    stop_data_redis &&
    stop_storage_redis
}

list() {
    cd redis_conf
    local auth_s=$(get_db_status $AUTH_REDIS_DUYOG_PID)
    local data_s=$(get_db_status $DATA_REDIS_DUYOG_PID)
    local stor_s=$(get_db_status $STORAGE_REDIS_DUYOG_PID)
    echo ""
    echo "+-------+----------+---------+"
    echo "| db    | server   | status  |"
    echo "+-------+----------+---------+"
    echo "|       | auth     | $auth_s |"
    echo "| redis | data     | $data_s |"
    echo "|       | storage  | $stor_s |"
    echo "+-------+----------+---------+"
    echo ""
}

get_db_status() {
     if [ -f  $1 ]; then
        echo "started"
    else
        echo "stopped"
    fi
}

case $1 in
    "redis")
        cd redis_conf
        case $2 in
            "start")
                case $3 in
                    "") start_all_redis;;
                    "auth") start_auth_redis;;
                    "data") start_data_redis;;
                    "storage") start_storage_redis;;
                esac;;
            
            "stop")
                case $3 in
                    "") stop_all_redis;;
                    "auth") stop_auth_redis;;
                    "data") stop_data_redis;;
                    "storage") stop_storage_redis;;
                esac;;
        esac;;

    "list")
        list;;
    
    "stop")
        cd redis_conf
        stop_all_redis;;
    
    "start")
        cd redis_conf
        start_all_redis;;
esac