import pytest
import requests
import json
from utils.generic_test_case import generate_test_cases
from utils.equivalence_classes import EquivalenceClasses as eq

test_cases_FETCH_1 = generate_test_cases(1,
    eq.communityId()
)

test_cases_FETCH_2 = generate_test_cases(1,
    eq.planId()
)

test_cases_FETCH_3 = generate_test_cases(2,
    eq.communityId(), eq.planId()
)
test_cases_POST_BULK_2 = generate_test_cases(4,
    eq.communityId(), eq.planId(), eq.communityId(), eq.planId(),
)

class TestCommunityPlan:
    base_url = "http://localhost:8098/"

    @pytest.mark.parametrize("CommunityId, expected",
                             test_cases_FETCH_1)
    def test_FETCH_1(self, CommunityId, expected):
        parameters = {"communityId": CommunityId}

        response = requests.get(self.base_url + "community-plan/",
                                params=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"

    @pytest.mark.parametrize("PlanId, expected",
                             test_cases_FETCH_2)
    def test_FETCH_2(self, PlanId, expected):
        parameters = {"planId": PlanId}

        response = requests.get(self.base_url + "community-plan/",
                                params=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} for PlanId={PlanId} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} for PlanId={PlanId} where data is {json.dumps(response.json(), indent=4)}"


    @pytest.mark.parametrize("CommunityId, PlanId, expected",
                             test_cases_FETCH_3)
    def test_FETCH_3(self, CommunityId, PlanId, expected):
        parameters = {"communityId": CommunityId,
                      "planId": PlanId}

        response = requests.get(self.base_url + "community-plan/",
                                params=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"

    @pytest.mark.parametrize("CommunityId, PlanId, expected",
                             test_cases_FETCH_3)
    def test_POST(self, CommunityId, PlanId, expected):
        parameters = {"communityId": CommunityId,
                      "planId": PlanId}

        response = requests.post(self.base_url + "community-plan/",
                                json=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"

    @pytest.mark.parametrize("parameters",
                             [{}, {"invalid": "format"},'{"malformed": json'])
    def test_POST_bulk_0(self,parameters):
        parameters = {}

        response = requests.post(self.base_url + "community-plan/",
                                json=parameters)

        if False:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"


    @pytest.mark.parametrize("CommunityId, PlanId, expected",
                             test_cases_FETCH_3)
    def test_POST_bulk_1(self, CommunityId, PlanId, expected):
        parameters = {"communityId": CommunityId,
                      "planId": PlanId}

        response = requests.post(self.base_url + "community-plan/",
                                json=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"

    @pytest.mark.parametrize("CommunityId0, PlanId0, CommunityId1, PlanId1, expected",
                             test_cases_POST_BULK_2)
    def test_POST_bulk_2(self, CommunityId0, PlanId0, CommunityId1,
                         PlanId1, expected):
        parameters = {"communityId": CommunityId0,
                      "planId": PlanId0,
                      "communityId": CommunityId1,
                      "planId": PlanId1}

        response = requests.post(self.base_url + "community-plan/",
                                json=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"

    @pytest.mark.parametrize("CommunityId, PlanId, expected",
                             test_cases_FETCH_3)
    def test_GET(self, CommunityId, PlanId, expected):
        parameters = {"communityId": CommunityId,
                      "planId": PlanId}

        response = requests.get(self.base_url + "community-plan/" +
                                CommunityId + "/" + PlanId + "/")

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"


    @pytest.mark.parametrize("CommunityId, PlanId, expected",
                             test_cases_FETCH_3)
    def test_DELETE(self, CommunityId, PlanId, expected):
        parameters = {"communityId": CommunityId,
                      "planId": PlanId}

        response = requests.delete(self.base_url + "community-plan/" +
                                CommunityId + "/" + PlanId + "/")

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} where data is {json.dumps(response.json(), indent=4)}"
