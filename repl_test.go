


func TestCleanInput(t *testing.T) {
    // ...
    // start by creating a slice of structs:

    cases := []struct{
        input string
        expected []string
    }{
            {
                input: "  hello world  ",
                expected: []string{"hello", "world"},
            },
            // add more cases here


    }

    // then loop over cases and run tests

    for _, c := range cases {
        actual := cleanInput(c.input)
        // check the actual length vs expected
        // if not a match, use t.Errorf to print msg and fail test

        for i := range actual {
            word := actual[i]
            expectedWord := c.expected[i]
            // check each word in slice
            // if no match, uses t.Errorf.....

            }
        }
    

        

    
}