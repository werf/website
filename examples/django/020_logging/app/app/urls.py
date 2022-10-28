from django.urls import include, path
from django.http import HttpResponseNotFound
urlpatterns = [
# [<snippet ping-route-import>]
    path('ping', include('ping.urls')),
# [<endsnippet ping-route-import>]
]

def handler404(request, *args, **argv):
    return HttpResponseNotFound()
