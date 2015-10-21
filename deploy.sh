printer(){
  printf '\n' && printf '=%.0s' {1..40} && printf '\n'
  echo $1
  printf '=%.0s' {1..40} && printf '\n'
}

export GOPATH=/srv/go
APP_DIR="$GOPATH/src/github.com/michigan-com/newsfetch"

cd $APP_DIR

printer "Updating newsfetch golang src ..."
git pull deploy live

printer "Download any required third part libraries ..."
go get -t ./...

printer "Removing old binary ..."
go clean -i

printer "Installing new binary ..."
make install

printer "Restarting newsfetch-toppages ..."
supervisorctl restart newsfetch-chartbeat

printer "Adding git release ..."
git tag -a $(cat VERSION) -m 'Production release'
git push --tags deploy master
