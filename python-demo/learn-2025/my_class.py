class people:
    def __init__(self, name, age):
        self.name = name
        self.age = age

    def speak(self):
        print("%s 说: 我 %d 岁，%d kg。" % (self.name, self.age, self.__weight))

    # 私有方法
    def __get_weight(self):
        return self.__weight

    name = "张三"
    age = 18
    # 私有属性
    __weight = 100


v = people("李四", 20)
v.speak()


class student(people):
    def __init__(self, name, age, grade):
        super().__init__(name, age)
        self.grade = grade

    # 覆写父类的方法
    def speak(self):
        print("%s 说: 我 %d 岁了，我在读 %d 年级" % (self.name, self.age, self.grade))

    grade = 1


v = student("王五", 22, 3)
v.speak()
