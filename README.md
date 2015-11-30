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

## Test
* `make test`

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

```
./newsfetch
```

## Usage

### Environment Variables

* `MONGO_URI` -- Save article, chartbeat, or recipe data to mongodb.
* `CHARTBEAT_API_KEY` -- API key for chartbeat, required for any `chartbeat` command.
* `DEBUG` -- Conditional debugging (DEBUG=\*, DEBUG=logger, DEBUG="logger,debugger").
* `GNAPI_DOMAIN` -- Domain for API endpoints, primarily used by the chartbeat command to send trigger to API that it got new data

### Fetch Articles

Grab a single article
```
$ ./newsfetch article -u "http://www.usatoday.com/story/news/nation/2015/10/15/baltimore-police-comissioner-protests/73973330/"
```

Grab all articles

```
$ ./newsfetch articles get -b
```

Specify specific site
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

### Chartbeat

Available APIs

### [Toppages](https://chartbeat.com/docs/api/explore/#endpoint=live/toppages/v3/+host=gizmodo.com)
Hit the Toppages API
```
$ ./newsfetch chartbeat toppages
```

### [Quickstats](https://chartbeat.com/docs/api/explore/#endpoint=live/quickstats/v4/+host=gizmodo.com)
Hit the quickstats API
```
$ ./newsfetch chartbeat quickstats
```

### [Topgeo](https://chartbeat.com/docs/api/explore/#endpoint=live/top_geo/v1/+host=gizmodo.com)
Hit the Topgeo API
```
$ ./newsfetch chartbeat topgeo
```

### All
Grab from all above APIs at once

```
$ ./newsfetch chartbeat all
```

Keep it running on loop every 30 seconds

```
$ ./newsfetch chartbeat all -l 30
```

### Copy Articles

Copy articles returned by Michigan API into the local Mongo database:

```
export MONGO_URI=mongodb://localhost:27017/mapi
newsfetch articles copy-from 'https://api.michigan.com/v1/news/freep/life?limit=1000'
```

### Body Extractor

```
$ ./newsfetch body -u [url]
```

Add title and make it the first line in the output

```
$ ./newsfetch body -t
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
