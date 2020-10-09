from django.shortcuts import render
from django.http import HttpResponse
from general.utils import generate_random_string
from general.utils import upload_file

# Create your views here.
def index(request):
    # Загружаем файл на s3
    upload_file("qwerty.txt", "test-container")

    return HttpResponse("Hello World! " + generate_random_string())