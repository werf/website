from django.shortcuts import render
from django.http import HttpResponse
from .models import Label
from django.shortcuts import redirect

# Create your views here.
def index(request):
    list_labels = Label.objects.all()

    return render(
        request,
        'index.html',
        context={'list':list_labels},
    )

def add_new(request):
    if request.method == "POST":
        getLabel = request.POST.get("label")
        Label.objects.create(label_text = getLabel)

        return redirect('/')

    return render(
        request,
        'add.html',
    )
