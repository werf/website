from django.http import HttpResponse
from django.shortcuts import render


def image(request):
    return HttpResponse(render(request, "./templates/image.html"))
