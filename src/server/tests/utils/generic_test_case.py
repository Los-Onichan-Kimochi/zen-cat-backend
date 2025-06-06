def create_sub_test_case(n_parameters, classes, indexes):
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

    if current_index < 0:
        return True, sub_test_case

    return False, sub_test_case

def create_test_case(test_case, n_parameters, equivalence_classes,
                     specific_cases):
    classes = []
    test_case_id = ""

    for parameter in range(n_parameters):
        class_type = (test_case >> parameter) & 1
        test_case_id += str(class_type)
        classes.append(equivalence_classes[parameter][class_type])

    indexes  = [0] * n_parameters

    while True:
        end, sub_test_case = create_sub_test_case(n_parameters, classes,
                                                  indexes)
        specific_cases .append(sub_test_case)

        if end:
            break

    print(test_case_id + "->" )

def generate_test_cases(n_parameters, *equivalence_classes):
    if len(equivalence_classes) != n_parameters:
        raise Exception("Invalid number of equivalence classes")

    total_test_cases = pow(2, n_parameters)
    specific_cases = []

    for test_case in range(total_test_cases):
        create_test_case(test_case, n_parameters, equivalence_classes,
                         specific_cases)

    return specific_cases
