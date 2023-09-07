-- name: InsertOne :one
insert into product (id, name, price, description, image)
values ($1, $2, $3, $4, $5)
returning id, name, price, description, image;