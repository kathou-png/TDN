import pytest

@pytest.fixture(scope="session")
def test_empty_db(client):
    """Start with a blank database."""
    rv = client.get('/')
    assert b'Jouons' in rv.data

def login(client, username):
    return client.post('/game', data=dict(
        user=username
    ), follow_redirects=True)

def test_login(client):
    rv = login(client, "abady")
    print(rv.data)
    assert b'Logout' in rv.data

def logout(client):
    return client.get('/logout', follow_redirects=True)

def test_logout(client):
    rv = logout(client)
    assert b'Un nom, s\'il vous plait' in rv.data