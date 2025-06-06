from utils.generic_test_case import generate_test_cases
from utils.equivalence_classes import EquivalenceClasses as eq

import requests
import json

test_cases_GET_fetch_3 = generate_test_cases(2,
    eq.communityId(), eq.planId()
)

c_1 ="communityId=e804b95a-a388-4751-b246-96fe97232d35"
p_1 = "planId=d1694efe-9a13-42d7-a9e8-4d629f9f2f35"

c_2 = "communityId=a1570014-f96c-4ba1-9ac6-e2aec2127910"
p_2 = "planId=6d222f80-8887-4cc2-b6a1-48d08cd2d742"


c_3 = "communityId=76035ca7-1d3b-4d7d-9091-fc55f7410e59"
p_3 = "planId=eb71f5e0-589d-4f1b-86e7-696c30e92bfe"


base = "http://localhost:8098/community-plan/plan//"

url = base + "?" + "" #+ "&" + p_1 #+ "&" + c_2 + "&" + p_2

response = requests.get(url)
data = response.json()

# Pretty print the JSON
#print(json.dumps(data, indent=4))
print(json.dumps(test_cases_GET_fetch_3, indent=4))
