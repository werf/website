from django.shortcuts import render
from django.http import HttpResponse

from general.utils import redis_counter_value_add

# Create your views here.
def index(request):
    return HttpResponse("views " + str(redis_counter_value_add()))