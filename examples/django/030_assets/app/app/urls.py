from django.urls import include, path

urlpatterns = [
    path('ping', include('ping.urls')),
# [<snippet image-route-import>]
    path('image', include('image.urls')),
# [<endsnippet image-route-import>]
]
