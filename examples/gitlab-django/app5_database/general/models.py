from django.db import models

# Create your models here.
class Label(models.Model):

    label_text = models.CharField(max_length=50)
