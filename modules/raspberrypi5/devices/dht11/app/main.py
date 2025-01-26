import os
import board
import adafruit_dht
from config import device, pin

os.environ['RASPBERRYPI_VERSION'] = device.getConfig("RASPBERRYPI_VERSION")

# Initial the dht device, with data pin connected to:
_device = getattr(adafruit_dht, device.type)
_board = getattr(board, pin)
dhtDevice = _device(_board)

try:
    d = dict();

    d['t'] = dhtDevice.temperature
    d['h'] = dhtDevice.humidity
    
    print(d)

except Exception as error:
    # dhtDevice.exit()
    # print(error)
    # return d
    pass
