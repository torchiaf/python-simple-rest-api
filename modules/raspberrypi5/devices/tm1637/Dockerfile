FROM debian AS build

ADD ./app /app

RUN apt-get update

# RUN apt remove python3-rpi.gpio
RUN apt-get install -y git python3-pip
RUN apt-get install -y python3-venv

RUN python3 -m venv pi-venv --system-site-packages
RUN pi-venv/bin/python3 -m pip install --upgrade pip setuptools wheel
RUN pi-venv/bin/pip install -U pyyaml pyinstaller rpi-lgpio

# TM1637
RUN git clone https://github.com/depklyon/raspberrypi-python-tm1637.git
RUN cd raspberrypi-python-tm1637 && /pi-venv/bin/python3 setup.py install

# Build init.py
RUN pi-venv/bin/pyinstaller --onefile --target-architecture arm64 -n tm1637 /app/main.py
