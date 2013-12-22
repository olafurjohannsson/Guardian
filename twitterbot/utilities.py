import re, time

class Utilities:

    @staticmethod
    def measure_time(func):
        now = time.time()
        func()
        return time.time() - now

    @staticmethod
    def validate_uri(uri):
        regex = re.compile(
        r'^(?:http|ftp)s?://' # http:// or https://
        r'(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+(?:[A-Z]{2,6}\.?|[A-Z0-9-]{2,}\.?)|' #domain...
        r'localhost|' #localhost...
        r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})' # ...or ip
        r'(?::\d+)?' # optional port
        r'(?:/?|[/?]\S+)$', re.IGNORECASE)
        try:
            r = regex.match(uri)
            if r.group():
                return True
        except:
            pass
        return False