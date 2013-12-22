from collections import namedtuple
import requests, pickle
from abc import ABCMeta, abstractmethod
from bs4 import BeautifulSoup
from utilities import *
from random import randint
from threading import Thread



# Tuple object
Quote = namedtuple('Quote', ['text', 'author'])


class QuoteException(Exception):
    def __init__(self, message):
        super().__init__(self, message)

# Base abstract class
class QuoteProvider:
    __metaclass__ = ABCMeta

    def __init__(self, filename='quotes.txt', url=None):
        self.url = url
        self.filename = filename
        self.quotes = list()



    # Public API

    def save(self, quote):
        ''' Saves a quote object in a pickle file '''
        try:
            with open(self.filename, 'ab') as quotefile:
                pickle.dump(quote, quotefile)
                return True
        except Exception as err:
            raise QuoteException("Could not save quote!\nErr: %s" % err)
        return False

    def exists(self, quote):
        ''' Checks if a quote object exists in a pickle file '''
        try:
            with open(self.filename, 'rb') as quotefile:
                while True:
                    data = pickle.load(quotefile)
                    if quote == data:
                        return True
        except:
            pass
        return False

    def randomize(self):
        ''' Return a random quote from the list '''
        if len(self.quotes) > 0:
            number = randint(0, len(self.quotes) - 1)
            return self.quotes[number]

    @abstractmethod
    def load(self):
        ''' Function that must be overwritten in sub-classes, it handles loading all the quotes into 'self.quotes' '''
        return


    # Private API

    @abstractmethod
    def __parse__(self, input):
        ''' Function that must be overwritten in sub-classes, it handles parsing the return output from 'self.html' '''
        return

    @abstractmethod
    def __fetch__(self, url):
        ''' abstract method that handles fetching data and adding to 'self.quotes' '''
        pass

    def __request__(self, url):
        ''' Make a GET request on a specific uri and return all the response from said GET request. '''
        url = url or self.url
        
        if not url or not Utilities.validate_uri(url):
            raise QuoteException("Url not valid!")
        
        r = requests.get(url)

        if r.status_code == 200:
            return r.text
        else:
            raise QuoteException("%s could not return quotes!" % self.url)

    def __html__(self, html):
        ''' Return a BeautifulSoup object from a given text string '''
        if not html:
            raise QuoteException("No html arg!")

        try:
            return BeautifulSoup(html)
        except Exception as err:
            raise QuoteException('Could not parse text into BeautifulSoup!')


# Subclass
class GoodreadQuote(QuoteProvider):
    def __init__(self):
        return super().__init__(url='')
    
    def __parse__(self, input):
        return 

    def load(self):
        return

    def __fetch__(self, url):
        return 

# Subclass
class BrainyQuote(QuoteProvider):
    def __init__(self):
        super().__init__(url='http://www.brainyquote.com/quotes/keywords/list%s.html')
    
    #  Overwritten
    def __parse__(self, input):
        try:
            if not input:
                raise QuoteException("Can't parse input!")
            # find all divs with correct class
            for div in [ x for x in input.find_all('div', attrs={'class': 'boxyPaddingBig'}) ]:
                # get text and author
                text, auth = [ y for y in div.text.split('\n') if y != '"' and y ]
                
                yield (text, auth)

        except Exception as err:
            raise QuoteException("Can't parse input!\nErr: %s" % err)
    
    def load(self):
        ''' Load all data in a multi threaded env '''
        threads = []
        for i in range(14): # 13 pages
            url = self.url % ('_{0}'.format(i) if i > 0 else '')
            t = Thread(target=self.__fetch__, args=(url,))
            threads.append(t)
            t.start()

        for thread in threads:
            thread.join()


    def __fetch__(self, url):
        ''' Utilizes all methods to fetch the data from pre specfied configuration '''
        # GET request for data
        data = self.__request__(url)
        
        # Change into HTML 
        html = self.__html__(data)

        # Parse html and iterate
        for data in self.__parse__(html):
            text, auth = data
            quote = Quote(text, auth)
            self.quotes.append(quote)


