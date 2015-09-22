CHANGELOG
=========

v.0.1.6 09-21-2015
------------------

* Chartbeat toppages command saves hourly max visits for documents in Article collection
* Passing around session variable instead of MongoUri

v.0.1.5 09-21-2015
------------------

* Removing old snapshots so that the most recent is the only one kept
* Added MONGOURI env variable to the tests for DB testing
* Created new custom conditional logger, pivoted based on DEBUG environment variable
* Removed the -v flag for all commands, please us DEBUG env variable to add logging statements

v.0.1.4 09-16-2015
------------------

* Added -l flag, to loop the command every n seconds. E.g (loop every 5 seconds): ./newsfetch chartbeat toppages -l 5
* Added supervisor conf

v0.1.3 09-16-2015
-----------------

* Forked text-summary to fix issues we are having with consistent summaries

v0.1.2 09-16-2015
------------------

* Added chartbeat argument. E.g. `./newsfetch chartbeat toppages`

v0.1.1 09-15-2015
------------------

* [FIX] Body extractor was throwing an unhandled exception when receiving an invalid url

v0.1.0 09-14-2015
-----------------

* New summarizer

v0.0.16 09-09-2015
------------------

* Determining a duplicate article is now based on article id and not article url.

v0.0.15 08-28-2015
------------------

* Added new summarizer based on github.com/neurosnap/sentences sentence tokenizer

v0.0.14 08-27-2015
------------------

* Created_at never gets updated

v0.0.13 08-21-2015
------------------

* Articles now being properly updated

v0.0.12 08-21-2015
------------------

* Revamped JSON unmarshalling for articles
* Fixed photo dimensions from not being saved properly
* Fixed duplicate article issue where id key was not being checked

v0.0.11 08-19-2015
------------------

* ExtractBodyFromUrl now requires a channel because it is being used as a goroutine
for concurrency
* Added LOGLEVEL environment variable for verbosity in cli output, it will override
and log level set in verbose mode
* Added extract test to ensure I didn't break the extractor

v0.0.10 08-19-2015
------------------

* Updated production servers to go 1.5
* Modified Makefile to use new linker variable format
* Updated format of timing output
* Updated verbosity output to > INFO, moved a lot of output to DEBUG

v0.0.9 08-19-2015
-----------------

* No longer removing all articles in mongo;
* Updating articles that are matched by their url

v0.0.8 08-19-2015
-----------------

* Moved commands to separate folder
* Formatting body content better

v0.0.1
------

* Init
