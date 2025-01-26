import builtins
import yaml
from pathlib import Path
from models import Module 

# Redefine print to show the output when running in containers
def print(*args):
    builtins.print(*args, sep=' ', end='\n', file=None, flush=True)

# Module configs
moduleDict = yaml.safe_load(Path('/sensors/module.yaml').read_text())
module = Module(moduleDict)

device = module.getDevice("tm1637")
clk = int(device.getConfig("TM1637_CLK"))
dio = int(device.getConfig("TM1637_DIO"))