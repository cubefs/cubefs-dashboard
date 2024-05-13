#! /bin/bash
set -e
BIN_PATH=bin
FRONT_PATH=frontend
BACK_PATH=backend
RootPath=$(pwd)

help() {
    cat <<EOF

Usage: ./build.sh [ -h | --help ] [ -a | --all ] [ -b | --backend ]
    -h, --help              show help info
    -a, --all     build backend and frontend
    -b, --backend             build backend
    -f, --frontend            build frontend
EOF
    exit 0
}

cmd="help"

ARGS=( "$@" )
for opt in ${ARGS[*]} ; do
    case "$opt" in
        -h|--help)
            help
            ;;
        -a|--all)
            cmd=build_all
            ;;
        -b|--backend)
            cmd=build_backend
            ;;
        -f|--frontend)
            cmd=build_frontend
            ;;
        *)
            ;;
    esac
done

Version=`git describe --abbrev=0 --tags 2>/dev/null`
BranchName=`git rev-parse --abbrev-ref HEAD 2>/dev/null`
CommitID=`git rev-parse HEAD 2>/dev/null`
BuildTime=`date +%Y-%m-%d\ %H:%M`

#SrcPath=${RootPath}/console
TargetFile=${2:-$RootPath/bin/cfs-gui}
TargetDic=$(dirname $TargetFile)
mkdir -p $TargetDic

LDFlags="-X github.com/cubefs/cubefs/proto.Version=${Version} \
    -X github.com/cubefs/cubefs/proto.CommitID=${CommitID} \
    -X github.com/cubefs/cubefs/proto.BranchName=${BranchName} \
    -X 'github.com/cubefs/cubefs/proto.BuildTime=${BuildTime}' "

function build_backend () {
  cd $BACK_PATH
  go build \
      -ldflags "${LDFlags}" \
      -o $TargetFile \
      *.go
  if [[ $? -ne 0 ]];then
    echo "go build error"
    exit 13
  fi
  cd -
  if [ -f "./$BACK_PATH/conf/config.yml" ]; then
      echo "config.yml exist, skipped"
  else
      cp -rf ./$BACK_PATH/conf/config.yml $TargetDic/config.yml
  fi
}

function build_frontend () {
  #compile frontend
  cd $FRONT_PATH
  npm install
  if [[ $? -ne 0 ]];then
    echo "npm install error"
    exit 11
  fi
  npm run build
  if [[ $? -ne 0 ]];then
    echo "npm run build error"
    exit 12
  fi
  cd -
  cp -rf ./$FRONT_PATH/dist $TargetDic/dist
}

function build_all() {
    build_backend
    build_frontend
}

case "-$cmd" in
    -help) help ;;
    -build_all) build_all ;;
    -build_backend) build_backend ;;
    -build_frontend) build_frontend ;;
    *) help ;;
esac