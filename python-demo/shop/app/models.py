from django.db import models


# Create your models here.
class User(models.Model):
    nick_name = models.CharField(max_length=200)
    mobile = models.CharField(max_length=200)
    class Meta:
        db_table = 'user'