import psycopg2
import random
from datetime import datetime, timedelta

DB_SETTINGS = {
    "dbname": "pepsi",
    "user": "pepsi",
    "password": "pepsi123",
    "host": "217.15.168.46",  # или "database", если внутри docker
    "port": 5432
}

