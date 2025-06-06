class EquivalenceClasses:
    def __init__(self):
        pass

    @staticmethod
    def id():
        valid_class = [
            "1",
            "2",
            "3"
        ]

        invalid_class = [
            "a",
            "b"
        ]

        return [valid_class, invalid_class]

    @staticmethod
    def email():
        valid_class = [
            "x",
            "y"
        ]

        invalid_class = [
            "99",
            "88"
        ]

        return [valid_class, invalid_class]

    @staticmethod
    def age():
        valid_class = [
            "H",
            "I"
        ]

        invalid_class = [
            "V"
        ]

        return [valid_class, invalid_class]
