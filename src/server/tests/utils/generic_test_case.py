def create_sub_test_case(n_parameters, classes, indexes, expected):
    sub_test_case = []

    for parameter in range(n_parameters):
        sub_test_case.append(classes[parameter][indexes[parameter]])

    current_index = n_parameters - 1

    while current_index >= 0:
        indexes[current_index] += 1

        if indexes[current_index] < len(classes[current_index]):
            break
        else:
            indexes[current_index] = 0
            current_index -= 1

    sub_test_case.append(expected)

    if current_index < 0:
        return True, sub_test_case

    return False, sub_test_case

def create_test_case(test_case, n_parameters, equivalence_classes,
                     specific_cases):
    classes = []
    expected = True
    idx = ""
    for parameter in range(n_parameters):
        class_type = (test_case >> parameter) & 1
        expected = expected and not bool(class_type)
        idx += str(class_type)
        classes.append(equivalence_classes[parameter][class_type])

    print(idx)


    indexes  = [0] * n_parameters

    while True:
        end, sub_test_case = create_sub_test_case(n_parameters, classes,
                                                  indexes, expected)
        print(sub_test_case)
        print(expected)
        specific_cases .append(sub_test_case)

        if end:

            break

def generate_test_cases(n_parameters, *equivalence_classes):
    if len(equivalence_classes) != n_parameters:
        raise Exception("Invalid number of equivalence classes")

    total_test_cases = pow(2, n_parameters)
    specific_cases = []

    for test_case in range(total_test_cases):
        create_test_case(test_case, n_parameters, equivalence_classes,
                         specific_cases)

    return specific_cases
