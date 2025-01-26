class DeviceConfig:
    def __init__(self, config):
        self.name = config["name"]
        self.value = config["value"]

class Device:
    def __init__(self, device):
        self.name = device["name"]
        self.type = device["type"]
        self.config = [DeviceConfig(x) for x in device["config"]]
        
    def getConfig(self, name):
        return [x.value for x in self.config if x.name == name][0]

class Module:
    def __init__(self, module):
        self.name = module["name"]
        self.type = module["type"]
        self.devices = [Device(x) for x in module["devices"]]
        
    def getDevice(self, name):
        return [x for x in self.devices if x.name == name][0]