
from threading import Thread
import threading


i = 1
locke = threading.Lock();

def someThreadFunction():
    global i
    for j in range(0,100000):
    	locke.acquire()
    	i = i - 1
    	locke.release()
# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")

def someThreadFunction2():
    global i
    for j in range(0,100000):
    	locke.acquire()
    	i = i + 1
    	locke.release()
# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")
    


def main():
    someThread = Thread(target = someThreadFunction, args = (),)
    someThread2 = Thread(target = someThreadFunction2, args = (),)
    someThread.start()
    someThread2.start()

    
    someThread.join()
    someThread2.join()

    global i

 
     
    print(i)




main()