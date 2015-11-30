# -*- coding: utf-8 -*-
import sys
import json

from pymongo import MongoClient
from bson import CodecOptions
from summarizer import summarize

class ArgumentError(Exception):
    pass

def connect(uri="mongodb://localhost:27017/mapi"):
    return MongoClient(uri)

def disconnect(client):
    return client.close()

def process_article_summaries(db, override=False):
    col = db.Article

    opts = CodecOptions(unicode_decode_error_handler='ignore')
    col = col.with_options(codec_options=opts)

    articles = None
    skipped = 0
    summarized = 0

    if override:
        articles = col.find()
    else:
        articles = col.find({
            "$or": [
                { "summary": { "$size": 0 } },
                { "summary": { "$exists": False } }
            ],
            "body": { "$ne": "" }
        })

        skipped = col.find({
            "summary": {
                "$not": { "$size": 0 },
                "$exists": True
            }
        }).count()

    for article in articles:
        #print("Processing {} ...".format(article['article_id']))
        summary = summarize(article['headline'], article['body'])
        col.update({ '_id': article['_id'] }, { '$set': { 'summary': summary } })
        summarized += 1

    return { 'summarized': summarized, 'skipped': skipped }

if __name__ == '__main__':
    if len(sys.argv) < 2:
        raise ArgumentError("Requires mongodb uri, eg: python summarize.py mongodb://localhost:27017/mapi")

    uri = sys.argv[1]

    override = False
    if len(sys.argv) > 2:
        o_arg = sys.argv[2].lower()
        if o_arg == "override" or o_arg == "1" or o_arg == "true":
            override = True

    client = connect(uri)
    db = client.get_default_database()

    results = process_article_summaries(db, override=override)

    disconnect(client)

    print(json.dumps(results))

