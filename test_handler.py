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
        {'body': data, 'Content-Type': 'application/sparql-update'},
        "blank_context")['statusCode'] == 200


def test_get_entities_unit():
    for test_case in test_cases:
        with open('fixtures/'+test_case['file'], 'r') as myfile:
            data = myfile.read()

        entities = handler.get_entities(data)

        for entity in test_case['out']:
            assert entity in entities
