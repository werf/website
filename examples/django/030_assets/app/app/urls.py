from django.urls import include, path
from django.http import HttpResponseNotFound

urlpatterns = [
    path('ping', include('ping.urls')),
# [<snippet image-route-import>]
    path('image', include('image.urls')),
# [<endsnippet image-route-import>]
]

def handler404(request, *args, **argv):
    return HttpResponseNotFound()
