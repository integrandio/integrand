# This is the integration test suite for LucidMQ
from Integrand import Integrand
import pytest
import string
import os
import random

INTEGRAND_URL = os.environ.get('INTEGRAND_URL', 'http://localhost:8000')  # The server's hostname or IP address
INTEGRAND_API_KEY = os.environ.get('INTEGRAND_API_KEY', '11111')

def get_random_string(length: int):
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    return result_str

class TestGlueAPI:
    # TODO: Setup by deleting all the topics
    # def test_get_all_glue_handlers_empty(self):
    #     integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
    #     response = integrand.GetAllGlueHandlers()
    #     print(response)
    #     assert response['status'] == 'success'
    #     assert response['data'] == []

    # def test_get_all_glue_handler_empty(self):
    #     route = get_random_string(5)
    #     integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
    #     with pytest.raises(Exception) as httpException:
    #         response = integrand.GetGlueHandler(route)
    #TODO: Check the error here
    #     print(httpException)
    
    def test_get_all_glue_handlers(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.GetAllGlueHandlers()
        print(response)

    def test_create_glue_handler(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.CreateGlueHandler(id, topicName)
        assert response['status'] == 'success'
        assert response['data']['id'] == id
        assert response['data']['topicName'] == topicName
    
    def test_get_glue_handler(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateGlueHandler(id, topicName)
        response = integrand.GetGlueHandler(id)
        assert response['status'] == 'success'
        assert response['data']['id'] == id
        assert response['data']['topicName'] == topicName
    
    def test_delete_glue_handler(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateGlueHandler(id, topicName)
        response = integrand.DeleteGlueHandler(id)
        assert response['status'] == 'success'

class TestTopicAPI:
    def test_get_all_topics(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response=integrand.GetAllTopics()
        print(response)
    
    def test_get_all_topics(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response=integrand.GetAllTopics()
        print(response)
