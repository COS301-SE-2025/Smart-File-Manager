import os

# Load secrets from env
try:
    SERVER_SECRET = os.environ['SFM_SERVER_SECRET']
except KeyError:
    SERVER_SECRET = None