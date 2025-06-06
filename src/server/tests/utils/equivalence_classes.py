class EquivalenceClasses:
    def __init__(self):
        pass

    @staticmethod
    def communityId():
        valid_class = [
            "e804b95a-a388-4751-b246-96fe97232d35"
            #"a1570014-f96c-4ba1-9ac6-e2aec2127910",
            #"76035ca7-1d3b-4d7d-9091-fc55f7410e59"
        ]

        invalid_class = [
            "e804b95aa3884751b24696fe97232d35",
            "e804b95a-a388-4751-b246-96fe97232d3",
            "e804b95a-a388-4751-b246-96fe97232d35X",
            "e8|4b@5a-#388-4751-b246-96fe97232d35"
        ]

        return [valid_class, invalid_class]



    @staticmethod
    def planId():
        valid_class = [
            "d1694efe-9a13-42d7-a9e8-4d629f9f2f35"
            #"6d222f80-8887-4cc2-b6a1-48d08cd2d742",
            #"eb71f5e0-589d-4f1b-86e7-696c30e92bfe"
        ]

        invalid_class = [
            "d1694efe9a1342d7a9e84d629f9f2f35",
            "e804b95a-a388-4751-b246-96fe97232d3",
            "e804b95a-a388-4751-b246-96fe97232d35X",
            "e8|4b@5a-#388-4751-b246-96fe97232d35"
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
