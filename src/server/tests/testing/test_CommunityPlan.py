import pytest


invalid_class_id = ["1", "2"]

valid_class_id = ["1", "2", "3"]


class TestCommunityPlan:
    ENDPOINT = "/community-plan/"

    @pytest.mark.parametrize("CommunityId", invalid_class_id)
    @pytest.mark.parametrize("PlanId", valid_class_id)
    def test_GET(self, CommunityId, PlanId, expected):
        assert True

    def test_POST_bulk(self, n, expected):
        assert True

    def test_DELETE_bulk(self, n, expected):
        assert True

    def test_GET_CommunityID_PlanID(self, n, expected):
        assert True

    def test_DELETE_CommunityID_PlanID(self, n, expected):
        assert True
