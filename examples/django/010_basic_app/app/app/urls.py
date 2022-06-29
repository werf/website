from django.urls import include, path

urlpatterns = [
# [<snippet ping-route-import>]
    path('ping', include('ping.urls')),
# [<endsnippet ping-route-import>]
]
