if [ -z "$1" ]; then
    echo "Please set environment name as argument" && exit 1
fi

if [ -z "$2" ]; then
    echo "Please provide commit short hash as argument" && exit 1
fi

function checkExit {
    exitcode=$?
    if [ $exitcode -ne 0 ]
    then
        echo "ERROR: script exited with code $exitcode"
        exit $exitcode
    fi
}

docker build -t warehouse-auth \
  --build-arg service=auth \
  --build-arg env=$1
checkExit
