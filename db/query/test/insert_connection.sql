-- name: InsertConnectionSeed :exec
INSERT INTO connection (id, node_id, "user", key, checksum)
VALUES (100, 100,
        'superuser',
        'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCfajDkRavEQiovz4NyyrjF7wToPFlt04I27YWL0C53nruTMt3bd5I8dqdamfCg+ebwj2TXucykMfo9UbSlUWvojVv00SODtG7iUTg7OUHu4womNhgzCC8iwAU7sllWw09ozNmvZZlsGWGda+4QrT3zx0x4XPwPDd2ejGegvJmv+bZGFln4azweiWEfohtdztjlw5MVVo4cbTwhsneyJJkcDW4snEKYafFGAbQF138I4/1sXyqWYDQpHGpGfN4t2WqEGDYV99L6x30Cb2ZNqbtYruztQaUW8q8qx1e8IYVKQzFtKR1Sh0eOxc8Qt4NXud4s00JxsUP76dzuD5aFZ5/l a.tikholoz@pc.local',
        '736acd8d9a40338142382e5ffb377179c19008cebe5557785a511f5e0c74ecec');
