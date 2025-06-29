// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: insert_connections.sql

package query

import (
	"context"
)

const insertConnectionsSeed = `-- name: InsertConnectionsSeed :exec
INSERT INTO connection (id, node_id, "user", encrypted_key, pub_key)
VALUES (100, 100,
        'superuser',
        '3d7656d5-65e0-4de5-b656-d565e07de57f',
        'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCfajDkRavEQiovz4NyyrjF7wToPFlt04I27YWL0C53nruTMt3bd5I8dqdamfCg+ebwj2TXucykMfo9UbSlUWvojVv00SODtG7iUTg7OUHu4womNhgzCC8iwAU7sllWw09ozNmvZZlsGWGda+4QrT3zx0x4XPwPDd2ejGegvJmv+bZGFln4azweiWEfohtdztjlw5MVVo4cbTwhsneyJJkcDW4snEKYafFGAbQF138I4/1sXyqWYDQpHGpGfN4t2WqEGDYV99L6x30Cb2ZNqbtYruztQaUW8q8qx1e8IYVKQzFtKR1Sh0eOxc8Qt4NXud4s00JxsUP76dzuD5aFZ5/l a.tikholoz@pc.local'),
       (101, 100,
        'happy-miner',
        'c426a698-c083-4e6e-a6a6-98c0837e6e3e',
        'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCfajDkRavEQiovz4NyyrjF7wToPFlt04I27YWL0C53nruTMt3bd5I8dqdamfCg+ebwj2TXucykMfo9UbSlUWvojVv00SODtG7iUTg7OUHu4womNhgzCC8iwAU7sllWw09ozNmvZZlsGWGda+4QrT3zx0x4XPwPDd2ejGegvJmv+bZGFln4azweiWEfohtdztjlw5MVVo4cbTwhsneyJJkcDW4snEKYafFGAbQF138I4/1sXyqWYDQpHGpGfN4t2WqEGDYV99L6x30Cb2ZNqbtYruztQaUW8q8qx1e8IYVKQzFtKR1Sh0eOxc8Qt4NXud4s00JxsUP76dzuD5aFZ5/l a.tikholoz@pc.local'),
       (102, 100,
        'superuser_2',
        '69f3a8d3-ec56-4919-b3a8-d3ec569919f7',
        'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCfajDkRavEQiovz4NyyrjF7wToPFlt04I27YWL0C53nruTMt3bd5I8dqdamfCg+ebwj2TXucykMfo9UbSlUWvojVv00SODtG7iUTg7OUHu4womNhgzCC8iwAU7sllWw09ozNmvZZlsGWGda+4QrT3zx0x4XPwPDd2ejGegvJmv+bZGFln4azweiWEfohtdztjlw5MVVo4cbTwhsneyJJkcDW4snEKYafFGAbQF138I4/1sXyqWYDQpHGpGfN4t2WqEGDYV99L6x30Cb2ZNqbtYruztQaUW8q8qx1e8IYVKQzFtKR1Sh0eOxc8Qt4NXud4s00JxsUP76dzuD5aFZ5/l a.tikholoz@pc.local')
`

func (q *Queries) InsertConnectionsSeed(ctx context.Context) error {
	_, err := q.db.Exec(ctx, insertConnectionsSeed)
	return err
}
