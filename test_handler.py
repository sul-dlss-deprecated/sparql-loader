import handler

test_cases = [
    {
        "name": "Test ETL organizations",
        "file": "etl_orgs.txt",
        "out": [
            'http://sul.stanford.edu/rialto/agents/orgs/vice-provost-for-student-affairs/dean-of-educational-resources/office-of-accessible-education/oae-operations',
            'http://sul.stanford.edu/rialto/agents/orgs/vice-provost-for-student-affairs/dean-of-educational-resources/womens-community-center'
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
    },
    {
        "name": "Test large sparql query",
        "file": "large_query.txt",
        "out": [
            'http://sul.stanford.edu/rialto/context/addresses/2b3d13e9-4cea-49a5-ac73-34c603461ab1_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/2b3d13e9-4cea-49a5-ac73-34c603461ab1',
            'http://sul.stanford.edu/rialto/publications/693df4b21c16dee82a08b966bb7070db',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_e4143964-9548-4959-b04d-e434fa4219d2',
            'http://sul.stanford.edu/rialto/context/addresses/e4143964-9548-4959-b04d-e434fa4219d2_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/e4143964-9548-4959-b04d-e434fa4219d2',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_6e99a7c0-eed4-4c4d-8f8e-ec4505c5a3a4',
            'http://sul.stanford.edu/rialto/context/addresses/6e99a7c0-eed4-4c4d-8f8e-ec4505c5a3a4_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/6e99a7c0-eed4-4c4d-8f8e-ec4505c5a3a4',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_494ed70c-ecad-42a1-8e14-7c947eaec4ab',
            'http://sul.stanford.edu/rialto/context/addresses/494ed70c-ecad-42a1-8e14-7c947eaec4ab_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/494ed70c-ecad-42a1-8e14-7c947eaec4ab',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_96418d53-d54d-4068-a33c-c2c875a41958',
            'http://sul.stanford.edu/rialto/context/addresses/96418d53-d54d-4068-a33c-c2c875a41958_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/96418d53-d54d-4068-a33c-c2c875a41958',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_3974dc70-83d8-49cb-8e5e-835a94822044',
            'http://sul.stanford.edu/rialto/context/addresses/3974dc70-83d8-49cb-8e5e-835a94822044_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/3974dc70-83d8-49cb-8e5e-835a94822044',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_7e577f2a-e2fd-4770-b961-c805ffcc5ad3',
            'http://sul.stanford.edu/rialto/context/addresses/7e577f2a-e2fd-4770-b961-c805ffcc5ad3_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/7e577f2a-e2fd-4770-b961-c805ffcc5ad3',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_49664e8c-d613-41db-b58e-b1492ed5c7bd',
            'http://sul.stanford.edu/rialto/context/addresses/49664e8c-d613-41db-b58e-b1492ed5c7bd_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/49664e8c-d613-41db-b58e-b1492ed5c7bd',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_35ae8118-8f72-4540-a595-caf8bc40d6ff',
            'http://sul.stanford.edu/rialto/context/addresses/35ae8118-8f72-4540-a595-caf8bc40d6ff_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/35ae8118-8f72-4540-a595-caf8bc40d6ff',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_3854c4ab-d91a-43b1-9e66-f2cf5e756a4c',
            'http://sul.stanford.edu/rialto/context/addresses/3854c4ab-d91a-43b1-9e66-f2cf5e756a4c_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/3854c4ab-d91a-43b1-9e66-f2cf5e756a4c',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_e4d7d217-59d3-4703-8ae9-cf525065a8fe',
            'http://sul.stanford.edu/rialto/context/addresses/e4d7d217-59d3-4703-8ae9-cf525065a8fe_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/e4d7d217-59d3-4703-8ae9-cf525065a8fe',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_eeb64a2d-dfe8-4a43-9cc2-8c5fa1698dfd',
            'http://sul.stanford.edu/rialto/context/addresses/eeb64a2d-dfe8-4a43-9cc2-8c5fa1698dfd_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/eeb64a2d-dfe8-4a43-9cc2-8c5fa1698dfd',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_2b319ae2-078d-468d-ae35-870ece493fbf',
            'http://sul.stanford.edu/rialto/context/addresses/2b319ae2-078d-468d-ae35-870ece493fbf_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/2b319ae2-078d-468d-ae35-870ece493fbf',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_3e50e3c8-e55d-49b6-814d-69cc27d1e475',
            'http://sul.stanford.edu/rialto/context/addresses/3e50e3c8-e55d-49b6-814d-69cc27d1e475_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/3e50e3c8-e55d-49b6-814d-69cc27d1e475',
            'http://sul.stanford.edu/rialto/context/relationships/WOS:000347715900024_6dd9fa77-bc2f-4cb6-bd93-9e3ef02d7e99',
            'http://sul.stanford.edu/rialto/context/addresses/6dd9fa77-bc2f-4cb6-bd93-9e3ef02d7e99_WOS:000347715900024',
            'http://sul.stanford.edu/rialto/agents/people/6dd9fa77-bc2f-4cb6-bd93-9e3ef02d7e99'
        ]
    }
]


def test_main_int():
    for test_case in test_cases:
        with open('fixtures/'+test_case['file'], 'r') as myfile:
            data = myfile.read()

    assert handler.main(
        {'body': data, 'content_type': 'application/sparql-update'},
        "blank_context")['statusCode'] == 200


def test_get_entities_unit():
    for test_case in test_cases:
        with open('fixtures/'+test_case['file'], 'r') as myfile:
            data = myfile.read()

        entities = handler.get_entities(data)

        for entity in test_case['out']:
            assert entity in entities
