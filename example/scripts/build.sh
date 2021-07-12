for dir in $(ls); do
  if [ -f ${dir}/main.go ]; then
    echo building $dir ...
    go build -o build/${dir} ${dir}/main.go
  fi
done
