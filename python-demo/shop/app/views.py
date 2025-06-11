from django.http import HttpResponse
from django.shortcuts import render

from . import models


# Create your views here.

def index(request):
    name = "张三"
    roles = ["管理员", "超级管理员", "普通管理员"]
    method = request.method
    return render(request, 'index.html', {"name": name, "roles": roles, "method": method})


def user_list(request):
    dataList = models.User.objects.all()
    return HttpResponse(dataList[0])
