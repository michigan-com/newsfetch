CHANGELOG
=========

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
