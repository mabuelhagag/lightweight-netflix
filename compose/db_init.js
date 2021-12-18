let res = [
    db.users.createIndex({ email: 1 }, { unique: true }),
    db.reviews.createIndex({ userId: 1 ,movieId: 1 }, { unique: true })
]

printjson(res)

