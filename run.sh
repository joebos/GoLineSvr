
go build

if [ $# -eq 0 ]
then
    echo "Missing text file path. Please enter: ./run.sh text_file_path"
else
    ./golinesvr "$1"
fi
