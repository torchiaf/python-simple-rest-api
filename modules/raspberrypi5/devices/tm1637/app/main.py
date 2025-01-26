import sys
import tm1637
from time import sleep
from config import print, clk, dio

tm = tm1637.TM1637(clk=clk, dio=dio)

try:
    arg = sys.argv[1]
    print('Received: {}'.format(arg))
    
    t = int(arg)
    tm.temperature(t) # show temperature 't*C'
except Exception as error:
    print(error)
    pass
# finally:
    # tm.write([0, 0, 0, 0])