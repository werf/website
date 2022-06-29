from django.http import HttpResponse, HttpResponseBadRequest
from .models import Talkers


def remember(request):
    req = request.GET

    if "answer" not in req:
        return HttpResponseBadRequest("You forgot the answer :(\n")

    if "name" not in req:
        return HttpResponseBadRequest("You forgot the name :(\n")

    data = Talkers()
    data.answer = req["answer"]
    data.name = req["name"]
    data.save()

    return HttpResponse("Got it.\n")


def say(request):
    obj = Talkers.objects.last()

    if obj is None:
        return HttpResponse("I have nothing to say.\n")

    return HttpResponse("{answer}, {name}!\n".format(answer=getattr(obj, "answer"), name=getattr(obj, "name")))
