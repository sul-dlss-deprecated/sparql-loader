import handler

def test_main():
    assert handler.main({'body': """INSERT DATA
{ <http://example/book3> <http://purl.org/dc/elements/1.1/title>    "A new book" ;
                         <http://purl.org/dc/elements/1.1/creator>  "A.N.Other" .
}""", 'content_type': 'application/sparql-update'}, "TEST")['ResponseMetadata']['HTTPStatusCode'] == 200
