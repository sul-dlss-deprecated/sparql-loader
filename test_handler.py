import handler

test_cases = [
    {
        "name": "Test ETL organizations",
        "file": "etl_orgs.txt",
        "out": [
            'http://sul.stanford.edu/rialto/agents/orgs/vice-provost-for-student-affairs/dean-of-educational-resources/office-of-accessible-education/oae-operations',  # noqa
            'http://sul.stanford.edu/rialto/agents/orgs/vice-provost-for-student-affairs/dean-of-educational-resources/womens-community-center'  # noqa
        ]
    },
    {
        "name": "Test sparql with graph",
        "file": "example_with_graph.txt",
        "out": [
            'http://sul.stanford.edu/rialto/context/names/75872',
            'http://sul.stanford.edu/rialto/context/positions/capFaculty_Bio-ABC_3784',
            'http://sul.stanford.edu/rialto/agents/orgs/Child_Health_Research_Institute',
            'http://sul.stanford.edu/rialto/agents/people/3784',
            'http://sul.stanford.edu/rialto/context/positions/capFaculty_Stanford_Neurosciences_Institute_3784'
        ]
    },
    {
        "name": "Test literal with quotes",
        "file": "quoted_literal.txt",
        "out": ['http://sul.stanford.edu/rialto/agents/orgs/b87e0d339b9997a81c7078fc3c227133']
    },
    {
        "name": "Test short name in literal",
        "file": "short_name.txt",
        "out": ['http://sul.stanford.edu/rialto/agents/people/189479']
    }
]


def test_main_int():
    for test_case in test_cases:
        with open('fixtures/'+test_case['file'], 'r') as myfile:
            data = myfile.read()

    assert handler.main(
        {'body': data, 'headers': {'Content-Type': 'application/sparql-update'}},
        "blank_context")['statusCode'] == 200


def test_main_unhappy_path_int():
    with open('fixtures/bad_insert.txt', 'r') as myfile:
        data = myfile.read()

    assert handler.main(
        {'body': data, 'headers': {'Content-Type': 'application/sparql-update'}},
        "blank_context")['statusCode'] == 400


def test_main_unappy_path_unit():
    with open('fixtures/decoded_query.txt', 'r') as myfile:
        data = myfile.read()

    assert handler.main({'body': data, 'headers': {'Content-Type': 'application/x-www-form-urlencoded'}},
                        "blank_contenxt") == {'body': '[MalformedRequest] query string not properly escaped',  # noqa
                                              'statusCode': 422}


def test_main_unknown_content_type_unit():
    with open('fixtures/encoded_query.txt', 'r') as myfile:
        data = myfile.read()

    assert handler.main({'body': data, 'headers': {'Content-Type': 'application/unknown'}},
                        "blank_contenxt") == {'body': "[MalformedRequest] Invalid Content-Type: 'application/unknown'",  # noqa
                                              'statusCode': 422}


def test_get_entities_unit():
    for test_case in test_cases:
        with open('fixtures/'+test_case['file'], 'r') as myfile:
            data = myfile.read()

        entities = handler.get_entities(data)

        for entity in test_case['out']:
            assert entity in entities


def test_not_malformed_query_unit():
    with open('fixtures/encoded_query.txt', 'r') as myfile:
        data = myfile.read()

    assert handler.is_malformed_query(data, "application/x-www-form-urlencoded") is None


def test_malformed_query_unit():
    with open('fixtures/decoded_query.txt', 'r') as myfile:
        data = myfile.read()

    assert handler.is_malformed_query(data, "application/x-www-form-urlencoded") == {
        'body': "[MalformedRequest] query string not properly escaped",
        'statusCode': 422}


def test_clean_content_type_unit():
    assert handler.clean_content_type('application/sparql-update') == 'application/sparql-update'
    assert handler.clean_content_type('application/sparql-update; charset=utf-8') == 'application/sparql-update'
