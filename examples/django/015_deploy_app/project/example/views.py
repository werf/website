from rest_framework import viewsets

from .models import Label
from .serizalizers import LabelSerializer


class LabelViewSet(viewsets.ModelViewSet):
    queryset = Label.objects.all().order_by('id')
    serializer_class = LabelSerializer
