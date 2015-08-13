printer(){
  printf '\n' && printf '=%.0s' {1..40} && printf '\n'
  echo $1
  printf '=%.0s' {1..40} && printf '\n'
}

APP_DIR="/srv/go/src/github.com/michigan-com/newsfetch"

cd $APP_DIR

printer "Updating newsfetch golang src ..."
go get -u github.com/michigan-com/newsfetch

printer "Installing new binary ..."
make install
