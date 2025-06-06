import pytest
import requests
import json
from utils.generic_test_case import generate_test_cases
from utils.equivalence_classes import EquivalenceClasses as eq

test_cases_GET_fetch_1 = generate_test_cases(1,
    eq.communityId()
)

test_cases_GET_fetch_2 = generate_test_cases(1,
    eq.planId()
)

test_cases_GET_fetch_3 = generate_test_cases(2,
    eq.communityId(), eq.planId()
)

class TestCommunityPlan:
    ENDPOINT = "http://localhost:8098/community-plan/"

    @pytest.mark.parametrize("CommunityId, expected", test_cases_GET_fetch_1)
    def test_GET_fetch_1(self, CommunityId, expected):
        parameters = {"communityId": CommunityId}

        response = requests.get(self.ENDPOINT, params=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} for CommunityId={CommunityId} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} for CommunityId={CommunityId} where data is {json.dumps(response.json(), indent=4)}"

    @pytest.mark.parametrize("PlanId, expected", test_cases_GET_fetch_2)
    def test_GET_fetch_2(self, PlanId, expected):
        parameters = {"planId": PlanId}

        response = requests.get(self.ENDPOINT, params=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} for PlanId={PlanId} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} for PlanId={PlanId} where data is {json.dumps(response.json(), indent=4)}"


    @pytest.mark.parametrize("CommunityId, PlanId, expected", test_cases_GET_fetch_3)
    def test_GET_fetch_3(self, CommunityId, PlanId, expected):
        parameters = {"communityId": CommunityId,
                      "planId": PlanId}

        response = requests.get(self.ENDPOINT, params=parameters)

        if expected:
            assert 200 <= response.status_code < 300, \
                f"Expected success but got {response.status_code} for CommunityId={CommunityId} where data is {json.dumps(response.json(), indent=4)}"
        else:
            assert 400 <= response.status_code < 600, \
                f"Expected failure but got {response.status_code} for CommunityId={CommunityId} where data is {json.dumps(response.json(), indent=4)}"


    def test_POST_bulk(self,):
        assert True

    def test_DELETE_bulk(self,):
        assert True

    def test_GET_CommunityID_PlanID(self,):
        assert True

    def test_DELETE_CommunityID_PlanID(self,):
        assert True
