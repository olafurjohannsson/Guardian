import traceback


class Logger:

    def __init__(self, message):
        self.log(message)

    def log(self, message):
        try:
            with open('log.txt', 'a+') as logfile:
                logfile.write(message + '\n')
        except:
            pass


