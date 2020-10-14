from django.shortcuts import render
from django.http import HttpResponse
from general.utils import batch_mailing

# Create your views here.
def index(request):
    return HttpResponse(batch_mailing("EMAIL"))