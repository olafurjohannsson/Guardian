from twitterbot import TwitterBot
from config import *
from logger import *
from utilities import *
from quoteproviders import QuoteException, BrainyQuote
from twitter import OAuth
from datetime import datetime
from time import time, sleep

def twitterify(input):
    return '#%s' % ''.join([ x for x in input.split(' ') ])


def main():
    try:
        now = time()
        # Construct auth and twitterbot object
        auth = OAuth(config['access_token'], config['access_secret'], config['consumer_key'], config['consumer_secret'])
        twitterBot = TwitterBot(auth)

        # create brainyquote object and load quotes
        quote = BrainyQuote()
        quote.load()

        iter = 0
        while True:
            # Safeguard, just stop running after 1000 iterations
            if iter >= 1000:
                break

            random_quote = quote.randomize()

            if not quote.exists(random_quote):
                text, auth = random_quote
                auth = twitterify(auth)
                
                if len(auth) + len(text) > 140:
                    text = text[0:140 - len(auth)]

                tweet = '%s %s' % (text, auth)

                if twitterBot.tweet(tweet):
                    quote.save(random_quote)
                    break
            iter = iter + 1

    except Exception as err:
        Logger("Unhandled exception!: %s" % err)
    except QuoteException as qerr:
        Logger("Quote exception in program: %s" % qerr)
    finally:
        message = 'Finished running TwitterBot, date: {0}, running time: {1} secs'.format(str(datetime.now()), time() - now)
        Logger(message)


if __name__ == '__main__':
    while True:
        main()
        sleep(60 * 60)
        print("Sleeping...")
