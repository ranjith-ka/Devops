# MongoDB

udemy: mongodb-the-complete-developers-guide

```bash
make install-mongo  # Download and Install for first time
make run-mongo
```

#### CURD:

create, update, read, delete

### TODO

Collections, documents
Find, insert, update, updateMany

Note:

16MB - Max size for Documents

### Projection

---

Serach required data, projection will do in mongo server.

Embedded documents/Nested Documents/Array documents

To print the flightdata collections,

`db.flightData.find().pretty()`

`db.fightData.deleteMany({})`

`db.flightData.find({"status.description": "on-time"}).pretty()`

### Structure documents

1. No schema approach, Mongoish

### Data types

1. Text, Bool, Number(int32/64,Decimal), ObejectId, ISODate(Timestamp), Embedded Documents, Arrays, BSON
   `db.stats()`

### Embedded vs Reference

1. Schemas and Relationships is the way to embedded and reference

    - one <-> one
    - one <-> many
    - many <-> many
    - many <-> one

    ```json
    {
        "Name": "Myname",
        "Age": 23,
        "Address": "This is address",
        "AvaiableTime": ObjectId("Timing_Referece_ID")  // Reference approach
    }
    ```

2. `$lookup [from: "Name", localField: "Age", foreignField: "_id", as: "create"]`

    Allows to merge/Embeded documents for relation

3. Validation (Schemas on Collection)

    - `db.CreateCollection()` && `db.runCommand({collMod: 'collectionname'})`

    - validation level (Error, Warning, info)

4. Running Mongo with available options in local

    - WiredTiger is default storage Engine
    - `db.ShutdownServer()` // To stop the Mongo service

5. `insert()` vs `insertOne()` vs `insertMany()`

## References

The MongoDB Limits: https://docs.mongodb.com/manual/reference/limits/

The MongoDB Data Types: https://docs.mongodb.com/manual/reference/bson-types/

More on Schema Validation: https://docs.mongodb.com/manual/core/schema-validation/
