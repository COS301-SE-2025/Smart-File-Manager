# API Documentation

The following file should serve as documentation for how to interact with the Database core engine.
Please note the following:

The C++ server only defines a single endpoint called QRY which has the following URL


````
http://to_be_determined/QRY/
````

Regardless of what type of query (i.e. GET, SET, DELETE, UPDATE) is being passed, the endpoint only specifies the following parameters

### Method

POST

### Parameters

Note that all parameters should be passed as strings

- `userID` :        A string identifying the user making the query (Note: This is important as it is how the engine determines which "Database" to access)


- `server_secret` : An alphanumeric key stored in the .env variables of the NodeJS server. Since user authentication is handled by the NodeJS server
                    this enables the C++ server to ensure that these requests are being made by registered users without having to store user info
                    such as usernames and passwords. 

- `query`:      This is a query as defined and understood by the MPDB query parser (i.e. GET, SET, DELETE, UPDATE) the remainder of this document will be
                spent elaborting on the syntax of these queries

### Response

- `success` : A boolean True or False value indicating if a query was executed successfully.
- `message` : An error message as to why the query could not execute [Blank if Success == True]
-  `data`   : A list of JSON objects as returned by the query, if applicable (GET Queries) 

#### Example Request

```json
{
  "userID": "Blackbeard",
  "server_secret": "as)#)(AFJS82)24121SKLFIU9322386656fj209KKJJDAAKJKII",
  "query" : "GET ship.name == 'The Black Pearl' "
}
```

#### Example Response for successful query

```json
{
    "success" : true,
    "message" : "",
    "data" :    [
                    {
                        "id" : 1,
                        "data" :    {
                                        "name" : "Jack Sparrow",
                                        "age" : 27,
                                        "hobbies" : ["Pillaging", "Drinking", "Cheating", "Tortuga"],
                                        "ship" : {
                                                        "name" : "The Black Pearl",
                                                        "type" : "Frigate",
                                                        "cannons" : 23
                                                    }
                                    }
                    },

                    {
                        "id" : 17,
                        "data" :    {
                                        "name" : "Hector Barbossa",
                                        "age" : 43,
                                        "hobbies" : ["Sleeping", "Getting Cursed", "Losing a leg"],
                                        "ship" : {
                                                        "name" : "The Black Pearl",
                                                        "type" : "Frigate",
                                                        "cannons" : 23
                                                    }
                                    }   
                                    
                    }
                ]
}
```

#### Example Response for unsucessful query
```json
{
    "success" : false,
    "message" : "The query key 'gold_balance' could not be found",
    "data" : []
}
```


## Types of Queries and their syntax

Note the following keywords used in the queries and their meaning

- `collection`  : Refers to the NOSQL equivalent of a SQL table i.e. a bunch of JSON objects grouped together
- `key`         : Refers to the NOSQL equivalent of a SQL column i.e. a word describing the nature of some data
- `operator`    : Refers to a list of mathematical operators used to compare a key to a value. Includes [<, >, <=, >=, ==, !=, *] Note that "*" refers to 
                  a RegExp
- `value`       : Refers to the value that the key must be compared to via the operator
- `sort`        : Refers to either ASC or DESC for ascending or descending sort respectively
- `limit`       : A numerical value representing the amount of JSON objects to return
- `offset`      : A numerical value reprsenting the amount of values to skip before starting to append objects to response  
- `data`        : Refers to a syntactically correct and well formed JSON object
- `update_key`  : Refers to the key that must be updated (See UPDATE query)
- `update_val`  : Refers to the value that all update keys must be set to (See UPDATE query)
- `delete_key`  : Refers to a key to be deleted in all matching queries (See DELETE query)

## GET

**Syntax:**
GET \<collection> \<key> \<operator> \<value> \<sort> \<limit> \<offset>

**Notes on query type:**
* Retrieves all JSON objects (in their entirety) that matches the query logic
* Keywords are not case-sensitive i.e. there is no difference between GET and gEt or get
* All query information such as the \<collection> and the \<value> is case-sensitive meaning that a search for collection "Pirates" might not return data 
  while a search for \<collection> "PIRATES" could
* Query clauses: [\<sort>, \<limit>, \<offset>] are completely optional
* Query clause key is optional if no \<operator> and \<value> was passed. (See examples below)
* Query clauses \<operator> and \<value> may not appear alone i.e. if an \<operator> is passed then so must a \<value> and vice-versa. Passing a
  \<operator> and \<value> also neccesitates passing a \<key>

**Example of valid queries as passed to API query parameter:**

* "GET PIRATES"
* "GET PIRATES ship.name"
* "GET PIRATES ship.type == 'Man of War'"
* "GET PIRATES" ship.type == 'Frigate' DESC"
*  "GET PIRATES" ship.type == 'Frigate' ASC 2"
*  "GET PIRATES" ship.type == 'Frigate' DESC 4 7"

## SET

**Syntax**
SET \<collection> \<data>

**Notes on query type**
* Inserts JSON object as specified by \<data> into the specified \<collection>
* \<collection> is case-sensitive, thus ensure you are inserting into the correct collection

**Example of valid queries as passed to API query parameter:**
* "SET PIRATES   '{
                        "name" : "Joshamee Gibs",
                        "rank" : "First Mate",
                        "age" : 53,
                        "address" : {
                                        "street" : "7-th Lane",
                                        "city" : "Tortuga"
                                    }
                    }'
    "

## UPDATE

**Syntax**
UPDATE \<collection> \<update_key> \<update_val> \<GET_QUERY>

**Notes on query type**
* Used to update a particular key to a new value for all objects matching \<GET_QUERY>
* This query will firstly find all JSON objects matching the \<GET_QUERY>. Note that this query's syntax is exactly as specified in GET
* After retrieving the objects matching the \<GET_QUERY> it will proceed to update those object's \<update_key> to \<update_val> where applicable
* Note that the following are valid options for \<update_val> [new_val, ++increment_number, --_decrement_number, **_factor]

**Example of valid queries as passed to API query parameter:**
* "UPDATE PIRATES rum_left '1300' GET PIRATES ship.name == 'Queen Anne's Revenge' "
* "UPDATE PIRATES days_since_last_bath '++1' GET PIRATES ship.name == 'The Flying Dutchmen'"


## DELETE

**Syntax**
DELETE \<delete_key> \<GET_QUERY>

**Notes on query type**
* Find all JSON objects matching the \<GET_QUERY> and deletes all \<delete_key> for those objects
* Valid options for \<delete_key> include [*, key_name] with "*" being all keys which implies the entire objects and key_name representing some key.

**Example of valid queries as passed to API query parameter:**
* "DELETE * GET PIRATES age >= 100"
* "DELETE ship.nickname GET ship.nickname == 'Small Boat'"