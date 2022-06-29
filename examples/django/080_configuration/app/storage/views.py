from django.http import HttpResponseBadRequest, HttpResponse
from django.shortcuts import render
from .forms import UploadFileForm
from .minio import upload_to_s3, download_from_s3


def upload(request):
    if request.method == 'POST':
        form = UploadFileForm(request.POST, request.FILES)
        if form.is_valid():
            upload_to_s3(request.FILES['file'].read())
            return HttpResponse('File uploaded.\n')
        else:
            HttpResponseBadRequest("You forgot to attach a file.\n")
    return HttpResponseBadRequest("The request method must be POST.")


def download(request):
    return HttpResponse(download_from_s3())
