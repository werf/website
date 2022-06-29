from django.urls import include, path
from ping import views as ping_view
from image import views as image_view
from talkers import views as talkers_view

urlpatterns = [
    path('ping', ping_view.ping, name='ping'),
    path('image', image_view.image, name='image'),
# [<snippet talkers-route-import>]
    path('remember', talkers_view.remember, name='remember'),
    path('say', talkers_view.say, name='say'),
# [<endsnippet talkers-route-import>]
]
