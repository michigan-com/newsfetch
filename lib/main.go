/* This package grabs data relating to Michigan.com's news sites.  Most of these components
are concerned with getting the necessary data from our sites, including chartbeat data.
After downloading the data from various web scraping, APIs, we then parse it and save
the data to a mongodb server.  This utility is essentially a data aggregator for
Michigan.com's API, mapi.
*/
package lib

var Sites = []string{
	"freep.com",
	"detroitnews.com",
	"battlecreekenquirer.com",
	"hometownlife.com",
	"lansingstatejournal.com",
	"livingstondaily.com",
	"thetimesherald.com",
}

var Sections = []string{
	"home",
	"news",
	"life",
	"sports",
	"business",
}
