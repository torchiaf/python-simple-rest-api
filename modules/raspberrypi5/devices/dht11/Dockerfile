FROM debian AS build

ADD ./app /app

RUN apt-get update
# RUN apt-get -y upgrade
RUN apt-get install -y git python3-pip
RUN apt-get install -y python3-venv
RUN apt-get install -y i2c-tools libgpiod-dev python3-libgpiod
RUN apt install --upgrade python3-setuptools

RUN python3 -m venv pi-venv --system-site-packages
RUN pi-venv/bin/python3 -m pip install --upgrade pip
RUN pi-venv/bin/pip install -U pyyaml pyinstaller

# DHT11
RUN pi-venv/bin/pip3 install --upgrade adafruit-python-shell
# RUN apt-get install -y wget
# RUN wget https://raw.githubusercontent.com/adafruit/Raspberry-Pi-Installer-Scripts/master/raspi-blinka.py
RUN pi-venv/bin/pip3 install --upgrade adafruit-blinka

RUN apt remove python3-rpi.gpio
RUN pi-venv/bin/pip3 uninstall -y RPi.GPIO
RUN pi-venv/bin/pip3 install --upgrade rpi-lgpio

RUN pi-venv/bin/pip3 install adafruit-circuitpython-dht
RUN apt-get install libgpiod2

# Build init.py
RUN pi-venv/bin/pyinstaller --onefile --target-architecture arm64 --add-data=pi-venv/lib64/python3.11/site-packages/adafruit_blinka:adafruit_blinka/ -n dht11 /app/main.py
