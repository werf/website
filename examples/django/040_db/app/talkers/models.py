from django.db import models


class Talkers(models.Model):
    answer = models.TextField()
    name = models.TextField()

    def __str__(self):
        return self.name
