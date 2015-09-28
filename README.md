# Michigan Newsfetch
Fetching the live feeds of all Gannett news sources in Michigan, parse them down to the essentials, and store them away.

# Setup
## Install
Set up your workspace as specified on the Golang documentation:
* [http://golang.org/doc/code.html#Workspaces]( http://golang.org/doc/code.html#Workspaces )

Use go to get the github repo (make sure you're in the `$GOPATH/src`)
* `go get github.com/michigan-com/newsfetch`


## Build
Go into the directory and build it
* `cd github.com/michigan-com/newsfetch`
* `make build`

## Version Bumping
Increment PATCH by 1
```
$ make bump
```

Specify version
```
$ make bump 0.1.0
```

## Run

Run the executable:

    export NEWSFETCH_MONGO_URL=mongodb://localhost:27017/mapi
    ./newsfetch

## Usage

### Body Extractor

```
$ ./newsfetch body -u [url]
```

Add title and make it the first line in the output

```
$ ./newsfetch body -t
```

### Fetch Articles

Grab all articles with body extractor

```
$ ./newsfetch articles get -b
```

Specify specific site no body extractor
```
$ ./newsfetch articles get -i freep.com
```

Specify specific site and specific section
```
$ ./newsfetch articles get -i freep.com -e sports
```

Specify multiple sites
```
$ ./newsfetch articles get -i freep.com,detroitnews.com -e sports
```

Just grab the article URL in the output
```
$ ./newsfetch articles get -i freep.com -e sports | awk -F"\t+" '{print $4}'
```

### Remove Articles

Removes all documents within Snapshot collection, which is what stores all the articles
```
$ ./newsfetch articles rm -s mongodb://localhost:27017/mapi
```

Skip prompt
```
$ ./newsfetch articles rm -ns mongodb://localhost:27017/mapi
```

### Copy Articles

Copy articles returned by Michigan API into the local Mongo database:

    export NEWSFETCH_MONGO_URL=mongodb://localhost:27017/mapi
    newsfetch articles copy-from 'https://api.michigan.com/v1/news/freep/life?limit=1000'

### Generate Summary

Creates a summary based on the title and the article body

```
$ ./newsfetch body -t | ./newsfetch summary
```

If you have the title, use the flag

```
$ ./newsfetch body | ./newsfetch summary -t "Cancer doc Farid Fata appeals 45-year prison sentence"
```

### Logging Output

All logging is determined by the `DEBUG` environment variable.

This will output all logging statements
```
DEBUG=* ./newsfetch articles get
```

This will output only the logger output
```
DEBUG=logger ./newsfetch articles get
```

This will output logger and chartbeat ouput
```
DEBUG=logger,chartbeat ./newsfetch chartbeat toppages -k APIKEY
```
