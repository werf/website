import random
import os
import io
from django.conf import settings
import redis

"""
This function return random ASCII-string
"""


def redis_counter_value_add():
    redis_instance = redis.StrictRedis(host=settings.REDIS_HOST,
                                  port=settings.REDIS_PORT, db=0)
    return redis_instance.incr('views')

    

