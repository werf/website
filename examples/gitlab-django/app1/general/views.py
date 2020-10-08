from django.shortcuts import render
from django.http import HttpResponse
from general.utils import generate_random_string

# Create your views here.
def index(request):
    return HttpResponse("Hello World! " + generate_random_string())
