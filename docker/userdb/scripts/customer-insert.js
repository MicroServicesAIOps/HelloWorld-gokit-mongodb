function get_results(result) {
    print(tojson(result));
}

function insert_customer(object) {
    print(db.customers.insert(object));
}

//pass eve
insert_customer({
    "_id": ObjectId("57a98d98e4b00679b4a830b2"),
    "username": "user",
    "password": "e2de7202bb2201842d041f6de201b10438369fb8",
    "salt": "6c1c6176e8b455ef37da13d953df971c249d0d8e",
});
//pass password
insert_customer({
    "_id": ObjectId("57a98d98e4b00679b4a830b5"),
    "username": "user1",
    "password": "8f31df4dcc25694aeb0c212118ae37bbd6e47bcd",
    "salt": "bd832b0e10c6882deabc5e8e60a37689e2b708c2",
});
print("_______CUSTOMER DATA_______");
db.customers.find().forEach(get_results);
print("______END CUSTOMER DATA_____");