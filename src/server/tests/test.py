from utils.generic_test_case import generate_test_cases
from utils.equivalence_classes import EquivalenceClasses as eq

XD = generate_test_cases(2,
                    eq.id(),
                    eq.email()
                    )

for i in XD:
    print(i)
