from logger import Logger
from config import *
from twitter import *

class TwitterBotException(object):

    def __init(self, message):
        self.message = message

    def __str__(self):
        return str(self.message)


class TwitterBot(object):

    '''
    Twitter bot that posts tweets on a given auth
    '''

    def __init__(self, auth):
        if not auth:
            raise TwitterBotException('No auth!')
        else:
            try:
                self.twitter = Twitter(auth=auth)

            except Exception as err:
                Logger(err)


    def tweet(self, tweet):
        if not tweet or not self.twitter:
            raise TwitterBotException('No tweet or self.twitter!')
        else:
            try:
                self.twitter.statuses.update(status=tweet)
                
                return True
            except Exception as err:
                Logger(err)

        return False
