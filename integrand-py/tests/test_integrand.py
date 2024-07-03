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

@pytest.fixture(scope="class")
def clean_up_topics():
    print("Setting Up Class")
    yield
    print("Tearing Down Class")
    integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
    response = integrand.GetAllTopics()
    topics = response['data']
    for topic in topics:
        integrand.DeleteTopic(topic['topicName'])

@pytest.mark.usefixtures("clean_up_topics")
class TestConnectorAPI:
    #TODO: Setup by deleting all the topics
    def test_get_all_connectors_empty(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.GetAllConnectors()
        assert response['status'] == 'success'
        assert len(response['data']) == 0

    def test_get_connector_empty(self):
        route = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        with pytest.raises(Exception) as httpException:
            response = integrand.GetConnector(route)
        #TODO: Check the error here
        print(httpException)

    def test_create_connector(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.CreateConnector(id, topicName)
        assert response['status'] == 'success'
        assert response['data']['id'] == id
        assert response['data']['topicName'] == topicName
        # Clean up
        integrand.DeleteConnector(id)
    
    def test_get_all_connectors(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateConnector(id, topicName)
        response = integrand.GetAllConnectors()
        assert response['status'] == 'success'
        assert len(response['data']) == 1
        # Clean up
        integrand.DeleteConnector(id)

    def test_get_connector(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateConnector(id, topicName)
        response = integrand.GetConnector(id)
        assert response['status'] == 'success'
        assert response['data']['id'] == id
        assert response['data']['topicName'] == topicName
        # Clean up
        integrand.DeleteConnector(id)
    
    def test_delete_connector(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateConnector(id, topicName)
        response = integrand.DeleteConnector(id)
        assert response['status'] == 'success'

class TestTopicAPI:
    def test_get_topics_empty(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.GetAllTopics()
        assert response['status'] == 'success'
        assert len(response['data']) == 0
    
    def test_get_topic_empty(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        with pytest.raises(Exception) as httpException:
            response = integrand.GetTopic(topicName)
        #TODO: Check the error here
        print(httpException)

    def test_create_topic(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.CreateTopic(topicName)
        assert response['status'] == 'success'
        assert response['data']['topicName'] == topicName
        assert response['data']['oldestOffset'] == 0
        assert response['data']['nextOffset'] == 0
        # Clean up
        integrand.DeleteTopic(topicName)

    def test_get_all_topics(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateTopic(topicName)
        response = integrand.GetAllTopics()
        assert response['status'] == 'success'
        assert len(response['data']) == 1
        # Clean up
        integrand.DeleteTopic(topicName)

    def test_get_topic(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateTopic(topicName)
        response = integrand.GetTopic(topicName)
        assert response['status'] == 'success'
        assert response['data']['topicName'] == topicName
        # Clean up
        integrand.DeleteTopic(topicName)
    
    def test_delete_topic(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateTopic(topicName)
        response = integrand.DeleteTopic(topicName)
        assert response['status'] == 'success'
    
class TestMessages():
    def test_send_message(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        createResponse = integrand.CreateConnector(id, topicName)
        data = {'hello': 'world'}
        res = integrand.EndpointRequest(id, createResponse['data']['securityKey'], data)
        assert res['status'] == 'success'
        response = integrand.GetEventsFromTopic(topicName, 0, 1)
        print(response)
        assert len(response['data']) == 1
        assert response['data'][0] == data
        # Clean up
        integrand.DeleteConnector(id)
        integrand.DeleteTopic(topicName)

